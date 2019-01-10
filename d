#/bin/sh

if [ -z ${DISPLAY+x} ]
then
    VIEWER='fbi'    
else
    VIEWER='sxiv'
fi

cd /tmp
dilbert
if [ $? -eq 0 ]
then
    $VIEWER dilbert-0.png
fi
