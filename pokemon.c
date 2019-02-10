//
//  main.c
//  PokemonC
//
//  Created by Julien on 22/09/2018.
//  Copyright © 2018 Julien Eyriès. All rights reserved.
//

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <time.h>

#define M 150
#define N 6

char *names[M];
uint64_t bits[M][3];
uint8_t answer[N];

void search() {
    for (int i= 0; i<N; i++) {
        answer[i] = 0;
    }

    int best = 0;

    for (int _i0 = 0; _i0 < M; _i0 ++) {
        int i0 = M-1 - _i0;

        printf("processing i0= %d\n", i0);
        fflush(stdout);
        clock_t start = clock();

        for (int i1 = 0; i1 < i0; i1++) {
            for (int i2 = 0; i2 < i1; i2 ++) {
                for (int i3 = 0; i3 < i2; i3 ++) {
                    for (int i4 = 0; i4 < i3; i4++) {
                        for (int i5 = 0; i5 < i4; i5++) {

                            uint64_t s0 = 0;
                            uint64_t s1 = 0;
                            uint64_t s2 = 0;

                            s0 |= bits[i0][0];
                            s1 |= bits[i0][1];
                            s2 |= bits[i0][2];

                            s0 |= bits[i1][0];
                            s1 |= bits[i1][1];
                            s2 |= bits[i1][2];

                            s0 |= bits[i2][0];
                            s1 |= bits[i2][1];
                            s2 |= bits[i2][2];

                            s0 |= bits[i3][0];
                            s1 |= bits[i3][1];
                            s2 |= bits[i3][2];

                            s0 |= bits[i4][0];
                            s1 |= bits[i4][1];
                            s2 |= bits[i4][2];

                            s0 |= bits[i5][0];
                            s1 |= bits[i5][1];
                            s2 |= bits[i5][2];

                            int score = __builtin_popcountll(s0) + __builtin_popcountll(s1) + __builtin_popcountll(s2);
                            if (score > best) {
                                best = score;
                                answer[0] = (uint8_t) i0;
                                answer[1] = (uint8_t) i1;
                                answer[2] = (uint8_t) i2;
                                answer[3] = (uint8_t) i3;
                                answer[4] = (uint8_t) i4;
                                answer[5] = (uint8_t) i5;
                            }
                        }
                    }
                }
            }
        }

        clock_t end = clock();
        double elapsed = (double)( end - start ) / (double)(CLOCKS_PER_SEC);
        printf("elapsed = %.3f s\n", elapsed);
        printf("best = %d\n", best);
        fflush(stdout);
    }

}

int main(int argc, const char * argv[]) {
    // insert code here...

    assert( argc == 2 );
    const char *path = argv[1];

    FILE *fp = fopen(path, "rt");

    int m = 0;
    char *line = NULL;
    size_t linecap = 0;
    ssize_t linelen;
    while ((linelen = getline(&line, &linecap, fp)) > 0) {
        //fwrite(line, linelen, 1, stdout);
        int i = 0;
        char *token;

        while ((token = strsep(&line, ",")) != NULL) {
            //printf("%s\n", token);
            if (i==0) {
                names[m] = strdup(token);
            } else {
                int x = atoi(token);
                bits[m][x / 64] |= 1LL << ( x % 64);
            }
            i += 1;
        }

        m += 1;
    }

    assert( m == M );
    search();

    for (int i= 0; i<N; i++) {
        puts(names[answer[i]]);
    }

    return 0;
}