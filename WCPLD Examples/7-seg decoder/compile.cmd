@echo off
set LIBCUPL=e:\Wincupl\Shared\atmel.dl
e:\Wincupl\Shared\wcupl.exe -jaxfsl -m2 %1
if errorlevel 1 (
echo errors found
.\wcupl\%1.lst
.\wcupl\%1.so
)