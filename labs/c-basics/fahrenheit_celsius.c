#include <stdio.h>


void convertFtoC(int f){
    printf("Fahrenheit: %3d, Celcius: %6.1f\n", f, (5.0/9.0)*(f-32));
}

int main(int argc, char **argv)
{
    int start,
        end,
        increment;
    if(argc>2){
        start   =   atoi(argv[1]);
        end   =   atoi(argv[2]);
        increment   =   atoi(argv[3]);
        for (start; start <= end; start = start + increment)
            convertFtoC(start);
    }
    else{
        if (argc    ==  0)
        {
            printf("Values were not introduced");
            return 0;
        }
        convertFtoC(atoi(argv[1]));

    }
    return 0;
}