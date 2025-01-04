![linux](https://github.com/nikonru/NiLang/actions/workflows/linux.yml/badge.svg)
![windows](https://github.com/nikonru/NiLang/actions/workflows/windows.yml/badge.svg)
# Why NiLang?
The part **NiLa** in the name stands for *Larry Niven* beloved author 
of the famous sci-fi book series *Ringworld*. 
The last two letters **ng** are reference to ISO 639 language code of *Ndonga* - one of the vibrant languages spoken in the Sub-Saharan Africa.

# What is this?
NiLang (Russian: НиЛанг) is a high level language for programming a bot from [TorLand](https://github.com/Slava2001/TorLand).
## File extension
We recommend use `.nil` as a file extension for files containing **NiLang** source code.
# Syntax
```
Using bot

Bool hungry = True
While hungry:
    Int maxEnergy = 1500
    ConsumeSunlight
    If GetEnergy > maxEnergy:
        Fork$ world::Forward
        hungry = False
Dir dir = GetDir
# you may remove `bot::`, since we have already written 'using bot'
bot::Move$ RotateClockwise$ dir  
```

```
Alias Numbers:
    One = 1
    Two = 2
    Four = 4

Scope names:
    max = 1000

Fun Bool F$max Int, default Bool:
    Using bot
    ConsumeSunlight
    If GetEnergy > max:
        Return True
    elif GetEnergy < max:
        Return False
    else:
        Return default

Fun Int G:
    Return 5
```
