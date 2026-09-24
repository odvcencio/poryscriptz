# Helpful patterns

These calls save a few common steps. They use the same script macros as direct calls. The compiler checks the argument shape and local label targets.

| Call | Use |
| --- | --- |
| `message(id)` | Open a message box. |
| `say(id)` | Open and close a message box. Add a wait when the player must read it first. |
| `trainer_battle_simple(trainer)` | Start a basic trainer battle. Check the result before a win event. |
| `give_item(item, count, bagFullLabel)` | Check bag space, then give an item. Define the full bag label. |
| `move_actor(actor, path)` | Start a movement path. |
| `move_actor_and_wait(actor, path)` | Start a movement path and wait for it. |
| `flag_set(flag)` | Set a flag. |
| `flag_clear(flag)` | Clear a flag. |
| `var_set(variable, value)` | Put a value in a variable. |
| `var_copy(destination, source)` | Copy one variable to another. |
| `warp_to(map, warp, x, y, direction)` | Warp to a map and position. |
| `if_flag_set(flag, label)` | Jump to a label when a flag is set. |
| `if_flag_unset(flag, label)` | Jump to a label when a flag is clear. |
| `if_var_eq(variable, value, label)` | Jump when a variable equals a value. |
| `if_var_ne(variable, value, label)` | Jump when a variable does not equal a value. |

A `movement` block gives a path its own label. End it with `EndMovement()`. For example:

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

Use `poryz commands` to find direct macros. The [cookbook](cookbook.md) shows complete files with these patterns.
