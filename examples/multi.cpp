#pragma once
#define FOO BAR
#define BAR BAZ
#define BAZ /*BAZ*/
/* test comment */  #define A /* inline comment */ 1 + 1/**/// single line comment FOO
/* multi
line

comment */FOO /* inline comment */ BAR BAZ // eol comment

#define BAZ 1

int main()
{
    return FOO + BAR + BAZ;
}

#undef BAR
#undef BAZ

#define FOO 2

int test()
{
    return FOO;
}

#define TEST test

int test2()
{
    return TEST();
}

#if FOO
int foo()
{
    return FOO;
}
#endif
