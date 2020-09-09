#include <stdio.h>

int main(int argc, char **argv) {
    if(argv[1]  ==  "-add"){
        char    *added  =mystradd(argv[2],argv[3]);
        printf("Initial Lenght      : %i\nNew String          : %s\nNew length          : %i",mystrlen(argv[2]),added,mystrlen(added));
    }
    else
    {
        if(argv[1]  ==  "-find"){
            int pos =   mystrfind(argv[2]);
            if(pos  !=  -1)
                printf("['%s'] string was found at [%i] position",argv[2],pos);
            else
            {
                printf("Substring not found");
            }
            
        }
    }
    
}
