# Cookbook

Each example is a full `.poryz` file. Replace message IDs and event flags with values from your own map. Run `poryz check examples/*.poryz` to check them. See the [command reference](command-reference.md) for every macro.

## Talk to an NPC

Lock the player, show one line, then give control back. The full file is [`examples/npc-dialogue.poryz`](../examples/npc-dialogue.poryz).

```poryz
package m

script Event {
    LockAll()
    FacePlayer()
    NPCMsg(0)
    WaitButton()
    CloseMsg()
    ReleaseAll()
    End()
}
```

## Show two pages

Use two message IDs when the talk needs a pause. The full file is [`examples/two-pages.poryz`](../examples/two-pages.poryz).

```poryz
package m

script Event {
    LockAll()
    NPCMsg(0)
    WaitButton()
    NPCMsg(1)
    WaitButton()
    CloseMsg()
    ReleaseAll()
    End()
}
```

## Read a sign

A sign can show a message without locking an NPC. The full file is [`examples/signpost.poryz`](../examples/signpost.poryz).

```poryz
package m

script Event {
    message(0)
    WaitButton()
    CloseMsg()
    End()
}
```

## Show a short message

The say helper shows and closes one message. The full file is [`examples/short-message.poryz`](../examples/short-message.poryz).

```poryz
package m

script Event {
    say(0)
    End()
}
```

## Give an item with a bag check

The helper sends the player to a full bag label when the item will not fit. The full file is [`examples/item-gift.poryz`](../examples/item-gift.poryz).

```poryz
package m

script Event {
    LockAll()
    give_item(ITEM_POTION, 1, _BagFull)
    ReleaseAll()
    End()
}
label _BagFull {
    CallStd(std_bag_is_full)
    ReleaseAll()
    End()
}
```

## Give an item without a bag check

Use this only when another step has checked the bag. The full file is [`examples/simple-gift.poryz`](../examples/simple-gift.poryz).

```poryz
package m

script Event {
    GiveItemNoCheck(ITEM_POTION, 1)
    End()
}
```

## Start a trainer battle

The helper starts a basic battle. Check the result before a win event. The full file is [`examples/trainer-battle.poryz`](../examples/trainer-battle.poryz).

```poryz
package m

script Event {
    LockAll()
    trainer_battle_simple(TRAINER_YOUNGSTER_JOEY)
    CheckBattleWon(VAR_SPECIAL_RESULT)
    ReleaseAll()
    End()
}
```

## Gate a trainer rematch

Use a flag to choose the first battle or a later message. The full file is [`examples/trainer-rematch.poryz`](../examples/trainer-rematch.poryz).

```poryz
package m

script Event {
    if flag(FLAG_UNK_042) {
        message(1)
        End()
    }
    trainer_battle_simple(TRAINER_YOUNGSTER_JOEY)
    flag_set(FLAG_UNK_042)
    End()
}
```

## Open an event with a flag

The else branch runs when the flag is clear. The full file is [`examples/flag-gate.poryz`](../examples/flag-gate.poryz).

```poryz
package m

script Event {
    if flag(FLAG_UNK_042) {
        message(1)
    } else {
        message(0)
    }
    End()
}
```

## Reset an event flag

Clear a flag when the event can run again. The full file is [`examples/flag-reset.poryz`](../examples/flag-reset.poryz).

```poryz
package m

script Event {
    flag_clear(FLAG_UNK_042)
    End()
}
```

## Hide a used item ball

Set the map flag after the item is given. The full file is [`examples/item-ball.poryz`](../examples/item-ball.poryz).

```poryz
package m

script Event {
    give_item(ITEM_POTION, 1, _BagFull)
    flag_set(FLAG_UNK_042)
    End()
}
label _BagFull {
    End()
}
```

## Set a script variable

Use a temp variable for a value that does not need to last. The full file is [`examples/set-variable.poryz`](../examples/set-variable.poryz).

```poryz
package m

script Event {
    var_set(VAR_TEMP_x4000, 1)
    End()
}
```

## Copy a script variable

Copy one value to another variable. The full file is [`examples/copy-variable.poryz`](../examples/copy-variable.poryz).

```poryz
package m

script Event {
    var_copy(VAR_TEMP_x4001, VAR_TEMP_x4000)
    End()
}
```

## Branch on a value

A branch target is a label in the same file. The full file is [`examples/variable-branch.poryz`](../examples/variable-branch.poryz).

```poryz
package m

script Event {
    if_var_eq(VAR_TEMP_x4000, 1, _Yes)
    message(0)
    End()
}
label _Yes {
    message(1)
    End()
}
```

## Ask yes or no

Store the answer, compare it, then jump to the yes label. The full file is [`examples/yes-no.poryz`](../examples/yes-no.poryz).

```poryz
package m

script Event {
    YesNo(VAR_SPECIAL_RESULT)
    Compare(VAR_SPECIAL_RESULT, 1)
    GoToIfEq(_Yes)
    message(0)
    End()
}
label _Yes {
    message(1)
    End()
}
```

## Handle more than one choice

A switch makes each result clear. The full file is [`examples/switch-choice.poryz`](../examples/switch-choice.poryz).

```poryz
package m

script Event {
    switch var(VAR_TEMP_x4000) {
    case 1:
        message(1)
    case 2:
        message(2)
    default:
        message(0)
    }
    End()
}
```

## Move an actor and wait

The next command runs after the path ends. The full file is [`examples/walk-and-wait.poryz`](../examples/walk-and-wait.poryz).

```poryz
package m

script Event {
    move_actor_and_wait(0, _Walk)
    End()
}
movement _Walk {
    WalkNormalNorth(2)
    EndMovement()
}
```

## Start movement and keep going

Use this when the next command can run during movement. The full file is [`examples/walk-in-background.poryz`](../examples/walk-in-background.poryz).

```poryz
package m

script Event {
    move_actor(0, _Walk)
    message(0)
    WaitMovement()
    End()
}
movement _Walk {
    WalkNormalEast(2)
    EndMovement()
}
```

## Face the player

Turn an NPC before a message. The full file is [`examples/face-player.poryz`](../examples/face-player.poryz).

```poryz
package m

script Event {
    LockAll()
    FacePlayer()
    message(0)
    WaitButton()
    CloseMsg()
    ReleaseAll()
    End()
}
```

## Play a sound

Play an existing sound before the message. The full file is [`examples/sound-effect.poryz`](../examples/sound-effect.poryz).

```poryz
package m

script Event {
    PlaySE(SEQ_SE_DP_SELECT)
    message(0)
    End()
}
```

## Warp to a map

Give the map, warp ID, coordinates, and direction. The full file is [`examples/warp.poryz`](../examples/warp.poryz).

```poryz
package m

script Event {
    warp_to(MAP_ROUTE_1, 0, 1, 1, DIR_SOUTH)
    End()
}
```

## Warp only after an event

Keep the player in place until the flag is set. The full file is [`examples/conditional-warp.poryz`](../examples/conditional-warp.poryz).

```poryz
package m

script Event {
    if flag(FLAG_UNK_042) {
        warp_to(MAP_ROUTE_1, 0, 1, 1, DIR_SOUTH)
    }
    End()
}
```

## Repeat while a flag is set

The loop stops after the body clears the flag. The full file is [`examples/repeat-until-clear.poryz`](../examples/repeat-until-clear.poryz).

```poryz
package m

script Event {
    for flag(FLAG_UNK_042) {
        message(0)
        flag_clear(FLAG_UNK_042)
    }
    End()
}
```

## Put the player name in a message

Buffer the name before a message that uses it. The full file is [`examples/player-name.poryz`](../examples/player-name.poryz).

```poryz
package m

script Event {
    BufferPlayersName(0)
    message(0)
    End()
}
```

## Give a party member

Check the result variable after the gift command. The full file is [`examples/gift-mon.poryz`](../examples/gift-mon.poryz).

```poryz
package m

script Event {
    GiveMon(SPECIES_PIKACHU, 5, ITEM_NONE, 0, 0, VAR_SPECIAL_RESULT)
    End()
}
```

## Call a helper label

A helper returns to the next command. The full file is [`examples/helper-label.poryz`](../examples/helper-label.poryz).

```poryz
package m

script Event {
    Call(_Helper)
    End()
}
label _Helper {
    PlaySE(SEQ_SE_DP_SELECT)
    Return()
}
```
