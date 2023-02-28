assume cs: code, ds: data

data segment
	dummy db 0Dh, 0Ah, '$'
	str1 db 100, 99 dup ('$')
	str2 db 100, 99 dup ('$')
data ends

code segment
include macro.asm
start:
	mov ax, data
	mov ds, ax
	mov ax, 0
	
	inputStr str1			; ввод строки в str1
	
	spaceDelete str1		; макрос для удаления пробелов
	
	printStr str2			; вывод строки
	
	exit				
	
code ends
end start