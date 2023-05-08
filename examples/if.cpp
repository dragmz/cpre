#if UNDEFINED
const int A = 1;
#endif

#define ZERO 0
#if ZERO
#if 1
const int B = 2;
#endif
#endif

#define NON_ZERO 1
#if NON_ZERO
const int C = 2;
#elif 1
const int D = 3;
#endif

#if 0
#elif 0
#else
const int E = 4;
#endif

#if 0
#elif 1
const int F = 5;
#endif

#define A C
#define B A
#define C B
#if C
const int G = 6;
#endif

#if 0
#include "test.h"
#pragma once
#endif