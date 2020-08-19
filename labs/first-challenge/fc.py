def getLength(array):
    l=0
    try:
        for e in array:
            l  =    l   + getLength(e)
    except:
        l   =   l   +   1
    return l

a=[]
print(getLength(a))
         
