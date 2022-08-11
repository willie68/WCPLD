@echo off
set LIBCUPL=e:\Wincupl\Shared\atmel.dl
e:\Wincupl\Shared\wcupl.exe -jaxfsl -m4 %1.wpld
rem c:\Wincupl\Shared\csim.exe -l g16v8 -u c:\Wincupl\Shared\Atmel.dl TEST
if errorlevel 1 (
.\wcupl\%1.lst
.\wcupl\%1.so
)