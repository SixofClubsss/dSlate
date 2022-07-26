Function InitializePrivate() Uint64
10 IF EXISTS("owner") == 0 THEN GOTO 30
20 RETURN 1 
30 STORE("owner", SIGNER())
40 RETURN 0
End Function

Function drawCardOne() Uint64  
10 DIM card1 as Uint64
20 LET card1 = 1+RANDOM(10)
30 IF EXISTS(card1) THEN GOTO 20
40 STORE(card1, card1)
50 MAPSTORE("mapCard1", card1)
60 RETURN 0
End Function

Function drawCardTwo() Uint64  
10 DIM card2 as Uint64
20 LET card2 = 1+RANDOM(10)
30 IF EXISTS(card2) THEN GOTO 20
40 STORE(card2, card2)
50 RETURN card2
End Function

Function DealCards() Uint64
10 drawCardOne()
20 STORE("Player 1 Card:", MAPGET("mapCard1"))
30 STORE("Player 2 Card:", drawCardTwo())
40 RETURN 0
End Function

Function ClearCards() Uint64
10 DIM i as Uint64
20 LET i = 0
30 DELETE(i)
40 LET i = i +1
50 IF i < 11 THEN GOTO 30
60 DELETE("Player 1 Card:")
70 DELETE("Player 2 Card:")
80 RETURN 0
End Function


