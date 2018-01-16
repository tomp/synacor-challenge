#!/usr/bin/env python3
#
import logging

# from scipy.special import comb
import gmpy2

logging.basicConfig(format="%(message)s", level=logging.INFO)
logger = logging.getLogger()


MODULUS = pow(2, 15)
# MODULUS = 0

class IterationError(Exception):
    pass

def r0_3_r1_n(n, r7):
    result = r7 * (r7 + n + 3) % MODULUS
    coef = gmpy2.mpz(n + 3)
    for i in range(2, n + 2):
        coef = gmpy2.divexact(coef * (n + 4 - i), i)
        result = (r7 * (result + coef)) % MODULUS
    result = (result + n) % MODULUS
    return result

def confirm(r0, r1, r7, maxstep=0, useshortcut=0, usecache=False):
    stack = []
    step = 0
    reg0, reg1, reg7 = r0, r1, r7

    cache = dict()
    seen_at = dict()

    def sub178b(lvl=0):
        nonlocal reg0, reg1, reg7, step, stack
        step += 1
        lvl += 1
        if maxstep and step > maxstep:
            raise IterationError("Too many steps!")

        if reg0 == 0:
            block = 0
        elif reg1 == 0:
            block = 1
        else:
            block = 2

        logger.debug("{:3d}: [{}]  r0:{:2d}  r1:{:2d}  s: {}".format(
            step, block, reg0, reg1, " ".join([str(v) for v in stack])))


        inputs = (reg0, reg1, reg7)
        if usecache:
            if not inputs in seen_at:
                seen_at[inputs] = len(stack)
            if inputs in cache:
                reg0, reg1, cached_steps = cache[inputs]
                # step += cached_steps
                return

        # shortcuts
        if reg0 <= useshortcut:
            if reg0 == 1:
                reg1 = reg7 + reg1
            elif reg0 == 2:
                reg1 = reg7 * (reg1 + 2) + reg1
            elif reg0 == 3:
                reg1 = r0_3_r1_n(reg1, reg7)
                # if reg1 == 0:
                #     reg1 = reg7 * (reg7 + 3)
                # elif reg1 == 1:
                #     reg1 = pow(reg7, 3) + 4 * pow(reg7, 2) + 6 * reg7 + 1
                # elif reg1 == 2:
                #     reg1 = (pow(reg7, 4) + 5 * pow(reg7, 3) +
                #             10 * pow(reg7, 2) + 10 * reg7 + 2)
                # else:
                #     reg1 = r0_3_r1_n(reg1, reg7)
            reg0 = reg1 + 1
            if MODULUS:
                reg0 = reg0 % MODULUS
                reg1 = reg1 % MODULUS
            return

        # original code
        if reg0 == 0:
            reg0 = reg1 + 1
        elif reg1 == 0:
            reg0 -= 1
            reg1 = reg7
            sub178b(lvl)
        else:
            stack.append(reg0)
            reg1 -= 1
            sub178b(lvl)
            reg1 = reg0
            reg0 = stack.pop() - 1
            sub178b(lvl)
        if MODULUS:
            reg0 = reg0 % MODULUS
            reg1 = reg1 % MODULUS

        if usecache and seen_at[inputs] == len(stack):
            cache[inputs] = (reg0, reg1, step)
            print("{} -> r0={}, r1={}".format(inputs, *cache[inputs]))
        return

    try:
        sub178b()
    except IterationError:
        pass

    return (reg0, reg1, step)

if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("-r0", type=int, default=4)
    parser.add_argument("-r1", type=int, default=1)
    parser.add_argument("-r7", type=int, default=1)
    parser.add_argument("-n", type=int, default=0)
    parser.add_argument("-d", "--debug", action='store_true')
    parser.add_argument("-s7", "--seq7")
    parser.add_argument("-s1", "--seq1")
    parser.add_argument("-x", "--shortcut", type=int, default=0)
    parser.add_argument("-y", "--cache", action='store_true')
    opt = parser.parse_args()

    assert opt.shortcut < 4

    logger.info("teleporter register confirmation")

    if opt.debug:
        logger.setLevel(logging.DEBUG)

    if opt.seq7:
        start, stop = [int(v) for v in opt.seq7.split(":")]
        for reg7 in range(start,stop+1):
            inputs = (opt.r0, opt.r1, reg7)
            (reg0, reg1, step) = confirm(opt.r0, opt.r1, reg7, opt.n,
                    useshortcut=opt.shortcut, usecache=opt.cache)
            logger.info("{:3d} steps,  {} --> r0:{:2d}  r1:{:2d}  r7:{:2d}".format(
                step, inputs, reg0, reg1, reg7))
    elif opt.seq1:
        start, stop = [int(v) for v in opt.seq1.split(":")]
        for reg1 in range(start, stop+1):
            inputs = (opt.r0, reg1, opt.r7)
            (reg0, reg1, step) = confirm(opt.r0, reg1, opt.r7, opt.n,
                    useshortcut=opt.shortcut, usecache=opt.cache)
            logger.info("{:3d} steps,  {} --> r0:{:2d}  r1:{:2d}  r7:{:2d}".format(
                step, inputs, reg0, reg1, opt.r7))
    else:
        inputs = (opt.r0, opt.r1, opt.r7)
        (reg0, reg1, step) = confirm(opt.r0, opt.r1, opt.r7, opt.n,
                useshortcut=opt.shortcut, usecache=opt.cache)
        logger.info("{:3d} steps,  {} --> r0:{:2d}  r1:{:2d}  r7:{:2d}".format(
            step, inputs, reg0, reg1, opt.r7))
