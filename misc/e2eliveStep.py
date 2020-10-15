#!/usr/bin/env python3
#
# For a private network, `goal network start` it, run indexer against
# that, verify that they get the same account data at the end.

import atexit
import glob
import gzip
import io
import json
import logging
import os
import random
import shutil
import sqlite3
import subprocess
import sys
import tempfile
import threading
import time
import urllib.request

from util import xrun, atexitrun, find_indexer, ensure_test_db, firstFromS3Prefix, deepeq

logger = logging.getLogger(__name__)

class e2elivestepper:
    def __init__(self):
        self.start = None
        self.args = None
        self.indexerurl = None
        self.healthurl = None
        self.accountsurl = None
        self.tempnet = None
        self.lastblock = None
        self.forwardRoundsAccounts = None
        self.indexer_bin = None
        self.psqlstring = None
        self.indexerout = None
        self.genesis = None

    def main(self):
        self.start = time.time()
        import argparse
        ap = argparse.ArgumentParser()
        ap.add_argument('--keep-temps', default=False, action='store_true')
        ap.add_argument('--indexer-bin', default=None, help='path to algorand-indexer binary, otherwise search PATH')
        ap.add_argument('--indexer-port', default=None, type=int, help='port to run indexer on. defaults to random in [4000,30000]')
        ap.add_argument('--connection-string', help='Use this connection string instead of attempting to manage a local database.')
        ap.add_argument('--source-net', help='Path to test network directory containing Primary and other nodes. May be a tar file.')
        ap.add_argument('--verbose', default=False, action='store_true')
        self.args = ap.parse_args()
        if self.args.verbose:
            logging.basicConfig(level=logging.DEBUG)
        else:
            logging.basicConfig(level=logging.INFO)

        try:
            self.setup_network()
            self.start_indexer()
            self.round_step_import()
            self.validate_endstate()
            self.validate_forward_reverse()
            dt = time.time() - self.start
            sys.stdout.write("indexer e2etest OK ({:.1f}s)\n".format(dt))
        except Exception as e:
            logger.error('error', exc_info=True)
            dt = time.time() - self.start
            sys.stdout.write("indexer e2etest FAILED ({:.1f}s)\n".format(dt))
            return 1
        return 0

    def setup_network(self):
        sourcenet = self.args.source_net
        source_is_tar = False
        if not sourcenet:
            e2edata = os.getenv('E2EDATA')
            sourcenet = e2edata and os.path.join(e2edata, 'net')
        if sourcenet and hassuffix(sourcenet, '.tar', '.tar.gz', '.tar.bz2', '.tar.xz'):
            source_is_tar = True
        tempdir = os.getenv('E2ETEMPDIR')
        if not tempdir:
            tempdir = tempfile.mkdtemp()
            logger.debug('created %r', tempdir)
            if not self.args.keep_temps:
                atexit.register(shutil.rmtree, tempdir, onerror=logger.error)
            else:
                atexit.register(print, "CLEANUP TODO\nrm -rf {!r}".format(tempdir))
        if not (source_is_tar or (sourcenet and os.path.isdir(sourcenet))):
            # fetch test data from S3
            bucket = 'algorand-testdata'
            import boto3
            from botocore.config import Config
            from botocore import UNSIGNED
            s3 = boto3.client('s3', config=Config(signature_version=UNSIGNED))
            tarname = 'net_done.tar.bz2'
            tarpath = os.path.join(tempdir, tarname)
            firstFromS3Prefix(s3, bucket, 'indexer/e2e2', tarname, outpath=tarpath)
            source_is_tar = True
            sourcenet = tarpath
        self.tempnet = os.path.join(tempdir, 'net')
        if source_is_tar:
            xrun(['tar', '-C', tempdir, '-x', '-f', sourcenet])
        else:
            xrun(['rsync', '-a', sourcenet + '/', self.tempnet + '/'])
        blockfiles = glob.glob(os.path.join(tempdir, 'net', 'Primary', '*', '*.block.sqlite'))
        self.lastblock = countblocks(blockfiles[0])
        xrun(['goal', 'network', 'start', '-r', self.tempnet])
        atexitrun(['goal', 'network', 'stop', '-r', self.tempnet])

    def read_genesis(self, algoddir):
        path = os.path.join(algoddir, 'genesis.json')
        with open(path, 'r') as fin:
            self.genesis = json.load(fin)

    def start_indexer(self):
        self.psqlstring = ensure_test_db(self.args.connection_string, self.args.keep_temps)
        algoddir = os.path.join(self.tempnet, 'Primary')
        self.read_genesis(algoddir)
        aiport = self.args.indexer_port or random.randint(4000,30000)
        self.indexer_bin = find_indexer(self.args.indexer_bin)
        cmd = [self.indexer_bin, 'daemon', '-P', self.psqlstring, '--dev-mode', '--server', ':{}'.format(aiport), '--no-algod']
        logger.debug("%s", ' '.join(map(repr,cmd)))
        indexerdp = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        self.indexerout = subslurp(indexerdp.stdout)
        self.indexerout.start()
        atexit.register(indexerdp.kill)
        self.indexerurl = 'http://localhost:{}/'.format(aiport)
        self.healthurl = self.indexerurl + 'health'
        self.accountsurl = self.indexerurl + 'v2/accounts'
        time.sleep(0.2)

    def round_step_import(self):
        self.forwardRoundsAccounts = {}
        errcount = 0
        ibenv = dict(os.environ)
        algoddir = os.path.join(self.tempnet, 'Primary')
        for xround in range(1,self.lastblock+1):
            ibenv['INDEXER_DEBUG_EXIT_ROUND'] = str(xround)
            xrun([self.indexer_bin, 'daemon', '-P', self.psqlstring, '--dev-mode', '--algod', algoddir, '--server', ''], env=ibenv)
            logger.debug('imported %d', xround)
            response = urllib.request.urlopen(self.accountsurl)
            if response.code != 200:
                raise Exception("{}: {} {!r}", self.accountsurl, response.code, response.read())
            raw = response.read()
            rob = json.loads(raw)
            for acct in rob['accounts']:
                if acct['round'] != xround:
                    logger.error('expected round %d but accout has %d', xround, acct['round'])
                    errcount += 1
                    if errcount > 10:
                        raise Exception('too many errors')
            self.forwardRoundsAccounts[xround] = raw

    def validate_endstate(self):
        for attempt in range(20):
            ok = tryhealthurl(self.healthurl, self.args.verbose, waitforround=self.lastblock)
            if ok:
                break
            time.sleep(0.5)
        if not ok:
            logger.error('could not get indexer health')
            sys.stderr.write(self.indexerout.dump())
            return 1
        try:
            algoddir = os.path.join(self.tempnet, 'Primary')
            xrun(['python3', 'misc/validate_accounting.py', '--verbose', '--algod', algoddir, '--indexer', self.indexerurl], timeout=20)
            xrun(['go', 'run', 'cmd/e2equeries/main.go', '-pg', self.psqlstring, '-q'], timeout=15)
        except Exception:
            sys.stderr.write(self.indexerout.dump())
            raise

    def validate_forward_reverse(self):
        errcount = 0
        special_addrs = set([self.genesis['fees'], self.genesis['rwd']])
        for xround in range(self.lastblock,0,-1):
            arurl = self.accountsurl + "?round={}".format(xround)
            logger.debug('GET %r', arurl)
            response = urllib.request.urlopen(arurl)
            if response.code != 200:
                raise Exception("{}: {} {!r}", arurl, response.code, response.read())
            logger.debug('reverse %d', xround)
            raw = response.read()
            oraw = self.forwardRoundsAccounts[xround]
            if raw == oraw:
                # easy, yay
                pass
            else:
                forward = json.loads(oraw)
                fbya = {x['address']:x for x in forward['accounts']}
                reverse = json.loads(raw)
                for ra in reverse['accounts']:
                    if ra.get('amount', 0) == 0:
                        continue
                    addr = ra['address']
                    if addr in special_addrs:
                        continue
                    fa = fbya.pop(addr, None)
                    if fa is None:
                        logger.error('round=%d reverse but not forward: %r', xround, ra)
                        errcount += 1
                        if errcount > 10:
                            raise Exception('too many errors')
                    else:
                        eqerr = []
                        if not deepeq(fa, ra, (), eqerr):
                            logger.error('round=%d neq forward=%r reverse=%r err=%r', xround, fa, ra, eqerr)
                            errcount += 1
                            if errcount > 10:
                                raise Exception('too many errors')

def hassuffix(x, *suffixes):
    for s in suffixes:
        if x.endswith(s):
            return True
    return False

def countblocks(path):
    db = sqlite3.connect(path)
    cursor = db.cursor()
    cursor.execute("SELECT max(rnd) FROM blocks")
    row = cursor.fetchone()
    cursor.close()
    db.close()
    return row[0]

def tryhealthurl(healthurl, verbose=False, waitforround=100):
    try:
        response = urllib.request.urlopen(healthurl)
        if response.code != 200:
            return False
        raw = response.read()
        logger.debug('health %r', raw)
        ob = json.loads(raw)
        rt = ob.get('message')
        if not rt:
            return False
        return int(rt) >= waitforround
    except Exception as e:
        if verbose:
            logging.warning('GET %s %s', healthurl, e, exc_info=True)
        return False

class subslurp:
    # asynchronously accumulate stdout or stderr from a subprocess and hold it for debugging if something goes wrong
    def __init__(self, f):
        self.f = f
        self.buf = io.BytesIO()
        self.gz = gzip.open(self.buf, 'wb')
        self.l = threading.Lock()
        self.t = None
    def run(self):
        for line in self.f:
            with self.l:
                if self.gz is None:
                    return
                self.gz.write(line)
    def dump(self):
        with self.l:
            self.gz.close()
            self.gz = None
        self.buf.seek(0)
        r = gzip.open(self.buf, 'rt')
        return r.read()
    def start(self):
        self.t = threading.Thread(target=self.run)
        self.t.daemon = True
        self.t.start()


if __name__ == '__main__':
    m = e2elivestepper()
    sys.exit(m.main())
