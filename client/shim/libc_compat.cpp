// https://reviews.llvm.org/D74712
// https://sourceware.org/git/?p=glibc.git;a=commit;h=7bdb921d70bf9f93948e2e311fef9ef439314e41

#include <math.h>

extern "C"
{
    __attribute__((__visibility__("default"), __cdecl__)) double __pow_finite(double a, double b)
    {
        return pow(a, b);
    }
}
