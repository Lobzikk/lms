@echo off
echo Welcome to setup batch script! If you want to skip some options, just press the Enter key!
IF not EXIST ServersInfo.txt (
    copy NUL ServersInfo.txt
)
IF not EXIST DBInfo.txt (
    copy NUL DBInfo.txt
)
set /p name= "Type as many unique names for the servers, as you can. Stop by inputting the stop word >>> "
:loop_names
    IF not %name% EQU stop (
        echo %name%>> ServersInfo.txt
        set /p name= ">>>"
        goto :loop_names
    )
set /p max= "Alright, now define how many workers could one server handle >>> "
echo %max%>> ServersInfo.txt
set /p url= "Now give me the url of your PostgreSQL DB  like that: <host:port> >>> "
echo %url%>> ServersInfo.txt
set /p username= "What's the name of user you are operating with in PostgreSQL DB? >>> "
echo %username%>> DBInfo.txt
set /p password= "Now give me the password (please!) >>> "
echo %password%>> DBInfo.txt
set /p dbname= "The last one step! What's the name of your database? (preferably 'LMS' or smth like that) >>> "
echo %dbname%>> DBInfo.txt
pause