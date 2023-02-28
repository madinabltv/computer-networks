inputStr macro strr
	lea dx, strr
	mov ah, 0Ah
	int 21h
	mov dx, offset dummy
	mov ah, 09h
	int 21h
endm

printStr macro st
	lea dx, st+2
	mov ah, 09h
	int 21h
endm

spaceDelete macro str
	lea si, str1+2
	lea di, str2+2
	proh:
		mov al, [si]
		inc si
		cmp al, 0Dh
		je endd
		cmp al, 32
		je proh
		mov [di], al
		inc di
		jmp proh
		endd:
endm

exit macro
	mov ax, 4C00h
	int 21h
endm