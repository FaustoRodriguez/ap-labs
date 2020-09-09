
int mystrlen(char *str){
    int i;
    for (i = 0; str[i] != '\0'; ++i);
    return i;
}

char *mystradd(char *origin, char *addition){
    int lenOrigin   =   mystrlen(origin),
        lenAddition =   mystrlen(addition),
        i;
    char    str[lenOrigin+lenAddition];
    for (i  =   0; i    <   lenOrigin; i++)
    {
        str[i]  =   origin[i];
    }
    for (; i    <   lenOrigin+lenAddition; i++)
    {
        str[i]  =   addition[i   -   lenOrigin];
    }
    return str;
}

int mystrfind(char *origin, char *substr){
    int lenOrigin   =   mystrlen(origin),
        lenSubstr   =   mystrlen(substr)-1,
        i,
        j;

    for(i   =   0;  i   +   lenSubstr   <   lenOrigin;  i++){
        for (j  =   0; j    <   lenSubstr   &&  origin[i  +   j]  ==   substr[j]; j++);
        if(j    ==   lenSubstr){
            return i;
        } 
    }

    return -1;
}
