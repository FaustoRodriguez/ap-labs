#include <stdio.h>
#include "logger.h"
#include <pthread.h>
#include <stdlib.h>
#define  NUM_THREADS 2000

long    **BUFFERS;
int     NUM_BUFFERS =   0;
char    *RESULT_MATRIX_FILE;
pthread_mutex_t *mtx;
pthread_t   thread[NUM_THREADS];

long    * readMatrix(char *filename),
        * getColumn(int col, long *matrix),
        * getRow(int row, long *matrix),
        * multiply(long *matA, long *matB),
        dotProduct(long *vecA, long *vecB);
int getLock(),
    releaseLock(int lock),
    saveResultMatrix(long *result);

struct arguments
{
    int row,
        col;
    long    *mat1,
            *mat2;
};


long mult(struct arguments * arg){

    int bufA = -1, bufB = -1;
    while(bufA == -1 || bufB == -1){
        if(bufA == -1){
            bufA = getLock();
        }
        if(bufB == -1){
            bufB = getLock();
        }
    }
    BUFFERS[bufA] = getRow(arg->row, arg->mat1);
    BUFFERS[bufB] = getColumn(arg->col, arg->mat1);
    long result = dotProduct(BUFFERS[bufA], BUFFERS[bufB]);
    free(BUFFERS[bufA]);
    free(BUFFERS[bufB]);
    free(arg);
    releaseLock(bufA);
    releaseLock(bufB);
    return result;
}


int main(int argc, char **argv){
    initLogger("stdout");

    if(argc < 5){
        errorf("Invalid number of arguments\n");
        errorf("Usage: -n <buffers> -out <filename>\n");
        return -1;
    }
    if(strcmp(argv[1],"-n") == 0 && strcmp(argv[3], "-out") == 0){
        NUM_BUFFERS = atoi(argv[2]);
        RESULT_MATRIX_FILE  = argv[4];
    }
    else{
        NUM_BUFFERS = atoi(argv[4]);
        RESULT_MATRIX_FILE  = argv[2];
    }
    BUFFERS = malloc(NUM_BUFFERS * sizeof(long *));
    mtx = malloc(NUM_BUFFERS * sizeof(pthread_mutex_t));
    for (int i = 0; i < NUM_BUFFERS; i++)
    {
        pthread_mutex_init(&mtx[i], NULL);
    }
    long    *matA = readMatrix("matA.dat"),
            *matB = readMatrix("matB.dat");
    long *result = multiply(matA,matB);
    infof("Finished\n");
    saveResultMatrix(result);
    free(matA);
    free(matB);
    free(mtx);
    free(BUFFERS);
    free(result);
    return 0;
}

long * readMatrix(char *filename){
    int size=0;
    FILE *f = fopen(filename, "r");
    if(f == NULL){
        errorf("Error, invalid Matrix datafile");
        exit(2);
    }
    char c;
    while ((c = fgetc(f)) != EOF)
    {
        if(c == '\n'){
            size++;
        }
    }
    rewind(f);
    long *matrix = malloc(size *sizeof(long));
    int i = 0;
    while (fscanf(f,"%ld", &matrix[i]) != EOF){
        i++;
    }
    fclose(f);
    infof("Matrix %s read",filename);
    return matrix;
}

long * getColumn(int col,long *matrix){
    long *column = malloc(NUM_THREADS * sizeof(long));
    for (int i = 0; i < NUM_THREADS; i++)
    {
        column[i]   = matrix[col];
        column      += NUM_THREADS;
    }
    return column;
}

long * getRow(int row, long *matrix){
    long *r = malloc(NUM_THREADS * sizeof(long));
    int n   = NUM_THREADS * (row);
    for (int i = 0; i < NUM_THREADS; i++)
    {
        r[i] = matrix[n+i];
    }
    return r;
}

int getLock(){
    for (int i = 0; i < NUM_BUFFERS; i++)
    {
        if (pthread_mutex_trylock(&mtx[i] == 0) == 0)
        {
            return i;
        }
    }
    return -1;
}

int releaseLock(int lock){
    return pthread_mutex_unlock(&mtx[lock]);
}

long dotProduct(long *vecA, long *vecB){
    long result = 0;
    for (int i = 0; i < NUM_THREADS; i++)
    {
        result += (vecA[i]*vecB[i]);
    }
    return result;
}

long * multiply(long *matA, long *matB){
    infof("Matrix multiplication starting");
    long *result = malloc(NUM_THREADS * NUM_THREADS * sizeof(long));
    for (int i = 0; i < NUM_THREADS; i++)
    {
        for (int j = 0; j < NUM_THREADS; j++)
        {
            struct arguments *arg = malloc(sizeof(struct arguments));
            arg-> row = i;
            arg-> col = j;
            arg-> mat1 = matA;
            arg-> mat2 = matB;
            pthread_create(&thread[j],NULL,(void * (*) (void *))mult,(void *)arg);
        }
        for (int j = 0; j < NUM_THREADS; j++)
        {
            void *res;
            pthread_join(thread[j], &res);
            result[NUM_THREADS * j + i] = (long) res;
        }
    }
    infof("Multiplication complete");
    return result;
}

int saveResultMatrix(long *result){
    FILE *f;
    f = fopen(RESULT_MATRIX_FILE,"w+");
    if (f == NULL)
    {
        return -1;
    }
    int size = NUM_THREADS * NUM_THREADS;
    for (int i = 0; i < size; i++)
    {
        fprintf(f,"%ld\n",result[i]);
    }
    fclose(f);
    return 0;    
}