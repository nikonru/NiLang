Using bot

Bool hungry = True
While hungry:
    int MaxEnergy = 1500
    ConsumeSunlight
    If GetEnergy > MaxEnergy:
        Fork$ world::Forward
        hungry = False
Dir dir = GetDir
# you may remove `bot::`, since we have already written 'using bot'
bot::Move$ RotateClockwise$ dir  