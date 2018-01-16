#include <stdint.h>
#include <stdio.h>

#define STACK_SIZE 1024 * 1024
#define TARGET 6

uint16_t reg0, reg1, reg7;
uint16_t stack[STACK_SIZE];
uint16_t sp = 0;

void calculate( uint16_t, uint16_t, uint16_t );
void sub178b( void );
void sub178b_v2( void );

int main(int argc, char *argv[])
{
    uint16_t r7;

    printf("teleporter confirmation code\n");

    for (r7=1; r7<32768; r7++) {
        uint16_t init0 = 4;
        uint16_t init1 = 1;
        uint16_t init7 = r7;

        calculate(init0, init1, init7);
        if (reg0 == TARGET) {
            printf("(%d, %d, %d) --> r0: %d  r1: %d\n",
                    init0, init1, init7, reg0, reg1);
        }
    }

}

void calculate(uint16_t r0, uint16_t r1, uint16_t r7)
{
    sp = 0;
    reg0 = r0;
    reg1 = r1;
    reg7 = r7;
    sub178b_v2();
}

void sub178b_v2( void )
{
    for (;;) {
        if (reg0 == 0) {
            break;
        }
        else if (reg0 == 1) {
            reg1 = (reg1 + reg7) & 0x7fff;
            break;
        }
        else if (reg0 == 2) {
            reg1 = (reg1 + (reg1 + 2) * reg7) & 0x7fff;
            break;
        }
        else if (reg1 == 0) {
            reg0--;
            reg1 = reg7;
        }
        else {
            stack[sp++] = reg0;
            reg1--;
            sub178b_v2();
            reg1 = reg0;
            reg0 = stack[--sp] - 1;
        } 
    } 
    reg0 = (reg1 + 1) & 0x7fff;
}

void sub178b( void )
{
    for (;;) {
        if (reg0 == 0) {
            reg0 = (reg1 + 1) & 0x7fff;
            return;
        }
        else if (reg1 == 0) {
            reg0--;
            reg1 = reg7;
        }
        else {
            stack[sp++] = reg0;
            reg1--;
            sub178b();
            reg1 = reg0;
            reg0 = stack[--sp] - 1;
        } 
    } 
}
