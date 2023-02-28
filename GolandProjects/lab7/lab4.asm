assume cs: code, ds: data

data segment
	dummy db 0Dh, 0Ah, '$'
	str1 db 100, 99 dup ('$')
	str2 db 100, 99 dup ('$')
	str3 db 100, 99 dup ('0')
	res db 100, 99 dup ('$')
	strerr db "error$"
data ends

code segment

free proc								; очистка
	xor ax, ax
	xor bx, bx
	xor cx, cx
	xor dx, dx
	xor si, si
	xor di, di
	ret
endp

print proc								; вывод на экран
	mov si, dx
	dec si
	lea di, res+2
	mov bh, 0
	p:
		mov bl, [si]
		cmp bl, 48
		jne nz
		cmp bh, 1
		je nz
		jmp point
		nz:
			mov bh, 1
			mov [di], bl
			inc di
			point:
				dec si
				mov al, [si]
				cmp al, 36
				jne p
				
	cmp bh, 0
	jne endprint
	add bh, 48
	mov [di], bh
	inc di
	
	endprint:
		mov bl, 36
		mov [di], bl
		xor al, al
		lea dx, res+2
		mov ah, 09h
		int 21h
	
	mov dx, offset dummy
	mov ah, 09h
	int 21h
	ret
endp

strinput proc						; ввод
	mov ah, 0Ah
	int 21h
	
	mov dx, offset dummy
	mov ah, 09h
	int 21h
	ret
strinput endp

convsub proc						; преобразование в число для вычитания
	sub ax, 48
	sub bx, 48
	
	cmp bh, 1
	jne noteq
	cmp al, 20
	jge erro
	jmp contin
	noteq:
		cmp al, 10
		jge erro
	contin:
		cmp bl, 10
		jge erro
		jmp return
	erro:
		mov ah, 9
		mov dx, offset strerr 
		int 21h
	enddd:
		mov ax, 4C00h
		int 21h
	return:
		ret
endp

convsum proc						; преобразование в число для сложения
	sub al, 48
	sub bl, 48
	
	cmp al, 10
	jge erro1
	cmp bl, 10
	jge erro1
	jmp return1
	erro1:
		mov ah, 9
		mov dx, offset strerr 
		int 21h
	enddd1:
		mov ax, 4C00h
		int 21h
	return1:
		ret
endp

decsum proc							; перевод в десятичную для сложения
	mov bl, 10
	mov si, dx
	div bl
	add ah, 48
	mov [si], ah
	xor ah, ah
	inc si
	inc dx
	cmp al, 0
	je nextt
	inc bh
	nextt:
		ret
decsum endp

decsub proc							; перевод в десятичную для вычитания
	mov si, dx
	add al, 48
	mov [si], al
	inc dx
	ret
endp

summ proc
	push bp
    mov bp, sp
	
	mov si, [bp+4]
	mov di, [bp+6]
	mov dl, [si]
	mov cl, [di]
	add si, dx
	add di, cx
	sub dx, cx
	push dx
	
	mov dx, [bp+8]
	inc dx
	l:
		mov al, [si]
		mov bl, [di]
		
		call convsum
		add al, bl
		cmp bh, 0
		je next
		inc ax
		xor bh, bh
		next:
			push si
			call decsum
			pop si
			dec si
			dec di
		loop l

	endproh:
		pop cx
		cmp cx, 0
		je endbit
		lp:
			mov al, [si]
			sub al, 48
			push si
			cmp bh, 0
			je next1
			inc ax
			xor bh, bh
			next1:
				call decsum
				pop si
				dec si
				loop lp
	endbit:
		cmp bh, 0
		je next2
		xor bh, bh
		push si
		call decsum
		pop si
		dec si
	next2:
		pop bp
		ret
endp

subb proc
	push bp
    mov bp, sp
	
	mov si, [bp+4]
	mov di, [bp+6]
	mov dl, [si]
	mov cl, [di]
	add si, dx
	add di, cx
	sub dx, cx
	push dx
	
	mov dx, [bp+8]
	inc dx
	passage:
		mov al, [si]
		mov bl, [di]
		
		cmp bh, 1
		jne nex
		dec al
		dec bh
		
		nex:
			cmp bl, al
			jle norm
			add al, 10
			inc bh
			call convsub
			norm:
				sub al, bl
				push si
				call decsub
				pop si
		dec si
		dec di
		loop passage
	endpas:
		pop cx
		cmp cx, 0
		je littlend
		ll:
			mov al, [si]
			
			cmp bh, 1
			jne nexx
			dec al
			dec bh
			
			nexx:
				cmp al, 47
				jg normm
				add al, 10
				inc bh
				normm:
					sub al, 48
					push si
					call decsub
					pop si
			dec si
			loop ll
	littlend:
		pop bp
		ret
endp

decmul proc
	mov bl, 10
	mov si, dx
	div bl
	add ah, 48
	mov [si], ah
	xor ah, ah
	inc dx
	cmp al, 0
	je nextt1
	mov bh, al
	nextt1:
		ret
endp

resmul proc
	mov si, dx
	mov bh, [si]
	sub bh, 48
	add al, bh
	mov bl, 10
	div bl
	add ah, 48
	mov [si], ah
	xor ah, ah
	mov bh, al
	inc dx
	ret
endp

mull proc
	push bp
    mov bp, sp
	
	mov si, [bp+4]
	mov di, [bp+6]
	mov cl, [si]
	mov dl, [di]
	
	cmp cl, 1
	je zero
	cmp dl, 1
	je zero
	jmp con
	
	zero:
		mov al, [si+1]
		cmp al, 48
		je zero1
		mov bl, [di+1]
		cmp bl, 48
		je zero1
		jmp con
		zero1:
			mov dx, [bp+8]
			inc dx
			mov al, 48
			mov si, dx
			mov [si], al
			inc dx
		pop bp
		ret
	
	con:
		mov ax, dx
		add si, cx
		add di, dx
		push si
		dec al
		
		mov dx, [bp+8]
		inc dx
		push dx
		push ax
		push cx
		
		beg:
			mov al, [si]
			mov bl, [di]
			
			call convsum
			mul bl
			cmp bh, 0
			je eqq
			add al, bh
			xor bh, bh
			eqq:
				push si
				call decmul
				pop si
			dec si
			loop beg
		
		mov si, dx
		add bh, 48
		mov [si], bh
		xor bh, bh
		inc dx
		
		pop bx
		pop ax
		push ax
		push bx
		cmp al, 0
		je end1
		
		firstCicle:
			pop cx
			pop ax
			
			pop dx
			inc dx
			
			pop si
			push si
			push dx
			push ax
			push cx
			
			dec di
			
			scndCicle:
				mov al, [si]
				mov bl, [di]
				
				call convsum
				mul bl
				add al, bh
				xor bh, bh
				push si
				call resmul
				pop si
				dec si
				loop scndCicle
				
			pop ax
			pop cx
			dec cx
			push cx
			inc cx
			push ax
			loop firstCicle
		
		mov si, dx
		cmp bh, 0
		je end1
		add bh, 48
		mov [si], bh
		inc dx
	end1:
		pop ax
		pop bx
		pop di
		pop si
	
	pop bp
	ret
	
endp

start:
	mov ax, data
	mov ds, ax
	mov ax, 0
	
	lea dx, str1
	call strinput
	
	lea dx, str2
	call strinput
	xor ax, ax
	
	lea si, str3+2
	mov [str3+2], 36
	push si
	lea si, str2+1
	push si
	lea si, str1+1
	push si
	
	call mull					; произведение беззнаковых чисел
	call print
	
	call free					; очистка регистров

	call subb					; разность беззнаковых чисел (первое число больше второго)
	call print

	call free
	
	call summ					; сумма беззнаковых (второе число не диннее первого)
	call print
	
	call free
	


	mov ax, 4C00h
	int 21h

code ends
end start