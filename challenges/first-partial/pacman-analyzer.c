#include <stdio.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <string.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>

void analizeLog(char *input, char *report);

struct Package
{
    char *name;
    char *installDate;
    char *lastUpdate;
    int updates;
    char *removalDate;
};

struct Segment
{
    char *date;
    char *name;
    char *status;
};


struct Package packages[1000];
int packagesTotal;

struct Segment getLineSegments(char *line)
{
    char    *date,
            *name,
            *status;
    struct Segment segment;
    int i   =   0,
        j   =   0;
    while(line[i]   !=  ']'){
        date[j++] =   line[++i];
    }
    i++;
    date    =   calloc(30,sizeof(char));
    name    =   calloc(30,sizeof(char));
    status  =   calloc(9,sizeof(char));
    while(line[i]   !=  ']'){
        i++;
    }
    i++;
    j   =   0;
    i++;
    while(line[i]   !=  ' '){
        status[j++]   =   line[i++];
    }
    i++;
    j   =   0;
    while(line[i]   !=  ' '){
        name[j++]   =   line[i++];
    }
    if(strcmp(status,"removed") ==  0 || strcmp(status,"upgraded")  ==  0 || strcmp(status,"installed") ==  0){
        segment.date    =   date;
        segment.name    =   name;
        segment.status  =   status;
        return segment;
    }
    return segment;
}

char *myGetLine(int file){
    int size = 200,
        c = 0,
        n;
    char *line;
    line    =   calloc(size, sizeof(char));
    while((n = read(file, line  +   c, size))   >   0){
        for(int i   =   c; i <   c   +   n; i++){
            if(line[i]  ==   '\n'){
                lseek(file, i   -   (c  +   n)  +   1, SEEK_CUR);
                line[i] =   '\0';
                return line;
            }
        }
        c   +=  n;
        line    =  realloc(line, c+size);
    }
    if(c    >   0){
        line[c] = '\0';
        return line;
    }
    return NULL;
}

int findPackage(char *name){
    for (int i = 0; i < packagesTotal; i++)
    {
        if (strcmp(name, packages[i].name) == 0)
        {
            return i;
        }
    }
    return -1;
}

int countPackages(char *status){
    int x=0;
    if(strcmp(status,"upgrades"))
    {
        for (int i = 0; i < packagesTotal; i++)
        {
            if (packages[i].lastUpdate  !=  NULL)
            {
                x++;
            }
        }
    }
    else{
        for (int i = 0; i < packagesTotal; i++)
        {
            if (packages[i].removalDate  !=  NULL)
            {
                x++;
            }
        }
    }
    return x;
}

int main(int argc, char **argv) {

    if (argc < 5) {
	printf("Usage:./pacman-analizer.o -input inputfile.txt -report reportfile.txt\n");
	return 1;
    }
    analizeLog(argv[2], argv[4]);
    return 0;
}

void analizeLog(char *input, char *report) {
    printf("Generating Report from: [%s] log file\n", input);
    int inputFile   =   open(input,O_RDONLY);
    char *line;
    struct Segment segment;
    line    =   myGetLine(inputFile);
    while (line !=  NULL)
    {
        segment =   getLineSegments(line);
        if(segment.name  !=  NULL){
            if(strcmp(segment.status,"installed")   ==  0){
                packages[packagesTotal].name    =   segment.name;
                packages[packagesTotal].installDate =   segment.date;
                packagesTotal++;
            }
            else
            {
                int packageID   =   findPackage(segment.name);
                if(packageID    >   -1){
                    if(strcmp(segment.status,"upgraded")    ==  0){
                        packages[packageID].lastUpdate  =   segment.date;
                        packages[packageID].updates++;
                    }
                    else{
                        if(strcmp(segment.status,"removed") ==  0){
                            packages[packageID].removalDate ==  segment.date;
                        }
                    }
                }
                else{
                    printf("Package [%s] not found",segment.name);
                }
            }   
        }
        line    =   myGetLine(inputFile);
    }
    close(inputFile);
    char n[5];
    int strSize,
        removedPackages =   countPackages("removed"),
        reportFile  =   open(report, O_WRONLY | O_CREAT | O_TRUNC, 0666);
    write(reportFile, "Pacman Packages Report\n", 23);
    write(reportFile, "----------------------\n", 23);
    write(reportFile, "- Installed packages : ", 23);
    strSize = sprintf(n, "%d\n", packagesTotal);
    write(reportFile, n, strSize);
    write(reportFile, "- Removed packages   : ", 23);
    strSize = sprintf(n, "%d\n", removedPackages);
    write(reportFile, n, strSize);
    write(reportFile, "- Upgraded packages  : ", 23);
    strSize = sprintf(n, "%d\n", countPackages("upgraded"));
    write(reportFile, n, strSize);
    write(reportFile, "- Current installed  : ", 23);
    strSize = sprintf(n, "%d\n", packagesTotal-removedPackages);
    write(reportFile, n, strSize);
    write(reportFile, "\nList of packages\n----------------\n", 35);
    for (int i = 0; i < packagesTotal; i++)
    {
        char *aux;
        write(reportFile, "- Package Name        : ", 24);
        write(reportFile, packages[i].name, strlen(packages[i].name));
        write(reportFile, "\n  - Install date      : ", 25);
        write(reportFile, packages[i].installDate, strlen(packages[i].installDate));
        write(reportFile, "\n  - Last update date  : ", 25);
        if (packages[i].lastUpdate  !=  NULL)
        {
            aux =   packages[i].lastUpdate;
        }
        else
        {
            aux =   "-";
        }
        write(reportFile, aux, strlen(aux));
        write(reportFile, "\n  - How many updates  : ", 25);
        strSize =   sprintf(n, "%d\n", packages[i].updates);
        write(reportFile, n, strSize);
        write(reportFile, "  - Removal date      : ", 24);
        if (packages[i].removalDate !=  NULL)
        {
            aux =   packages[i].removalDate;
        }
        else
        {
            aux =   "-";
        }
        write(reportFile, aux, strlen(aux));
        write(reportFile, "\n", 1);
    }
    close(reportFile);
    printf("Report is generated at: [%s]\n", report);
}