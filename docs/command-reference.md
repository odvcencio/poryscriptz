# Command reference

Generated from the public decomp's `script.inc` and `movement.inc` macro declarations. Call these names in a `.poryz` script or movement block. Optional parameters show their default. Use `poryz constants PREFIX` to find flag, variable, item, species, move, and map names.

| Command | Parameters | Area |
| --- | --- | --- |
| `ScrDef` | `offset` | Script |
| `ScrDefEnd` | `` | Script |
| `Noop` | `` | Script |
| `Dummy` | `` | Script |
| `End` | `` | Script |
| `Wait` | `frames, var` | Script |
| `LoadByte` | `reg, val` | Script |
| `LoadWord` | `reg, val` | Script |
| `LoadByteFromAddr` | `reg, addr` | Script |
| `WriteByteToAddr` | `addr, reg` | Script |
| `SetPtrByte` | `addr, val` | Script |
| `CopyLocal` | `to, from` | Script |
| `CopyByte` | `to, from` | Script |
| `CompareLocalToLocal` | `a, b` | Script |
| `CompareLocalToValue` | `reg, val` | Script |
| `CompareLocalToAddr` | `reg, addr` | Script |
| `CompareAddrToLocal` | `addr, reg` | Script |
| `CompareAddrToValue` | `addr, val` | Script |
| `CompareAddrToAddr` | `a, b` | Script |
| `CompareVarToValue` | `var, val` | Script |
| `CompareVarToVar` | `a, b` | Script |
| `RunScript` | `id` | Script |
| `CallStd` | `id` | Script |
| `RestartCurrentScript` | `` | Script |
| `GoTo` | `dest` | Script |
| `ObjectGoTo` | `object, dest` | Script |
| `BGGoTo` | `bg, dest` | Script |
| `DirectionGoTo` | `dir, dest` | Script |
| `Call` | `sub` | Script |
| `Return` | `` | Script |
| `GoToIf` | `condition, dest` | Script |
| `CallIf` | `condition, sub` | Script |
| `SetFlag` | `flag` | Script |
| `ClearFlag` | `flag` | Script |
| `CheckFlag` | `flag` | Script |
| `SetFlagVar` | `var` | Script |
| `ClearFlagVar` | `var` | Script |
| `CheckFlagVar` | `var_flag, var_dest` | Script |
| `SetTrainerFlag` | `var_or_trno` | Script |
| `ClearTrainerFlag` | `var_or_trno` | Script |
| `CheckTrainerFlag` | `var_or_trno` | Script |
| `AddVar` | `var, var_or_addend` | Script |
| `SubVar` | `var, var_or_addend` | Script |
| `SetVar` | `var, val` | Script |
| `CopyVar` | `dst, src` | Script |
| `SetOrCopyVar` | `dst, src` | Script |
| `NonNPCMsg` | `msg_id` | Script |
| `NPCMsg` | `msg_id` | Script |
| `NonNPCMsgVar` | `msg_id` | Script |
| `NPCMsgVar` | `arg0` | Script |
| `ScrCmd_048` | `arg0` | Script |
| `WaitABPress` | `` | Script |
| `WaitButton` | `` | Script |
| `WaitButtonOrDpad` | `` | Script |
| `OpenMsg` | `` | Script |
| `CloseMsg` | `` | Script |
| `HoldMsg` | `` | Script |
| `DirectionSignpost` | `message, type, map, out` | Script |
| `SetSignpostMap` | `type, map` | Script |
| `SetSignpostAction` | `cmd` | Script |
| `WaitSignpostAction` | `` | Script |
| `TrainerTips` | `message, out` | Script |
| `TrainerTipsEx` | `type, message` | Script |
| `WaitSignpost` | `out` | Script |
| `DirectionSignpostEx` | `type, map, message` | Script |
| `ScrCmd_061` | `` | Script |
| `ScrCmd_062` | `arg0, arg1, arg2, arg3, arg4, arg5` | Script |
| `YesNo` | `arg0` | Script |
| `ScrCmd_064` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_065` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_066` | `arg0, arg1` | Script |
| `ScrCmd_067` | `` | Script |
| `ScrCmd_068` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_069` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_070` | `arg0, arg1, arg2` | Script |
| `ScrCmd_071` | `` | Script |
| `ScrCmd_072` | `arg0` | Script |
| `PlaySE` | `arg0` | Script |
| `StopSE` | `arg0` | Script |
| `WaitSE` | `arg0` | Script |
| `PlayCry` | `arg0, arg1` | Script |
| `WaitCry` | `` | Script |
| `PlayFanfare` | `fanfare` | Script |
| `WaitFanfare` | `` | Script |
| `PlayBGM` | `bgm` | Script |
| `StopBGM` | `arg0` | Script |
| `ResetBGM` | `` | Script |
| `ScrCmd_083` | `arg0` | Script |
| `FadeOutBGM` | `target, frames` | Script |
| `FadeInBGM` | `arg0` | Script |
| `ScrCmd_086` | `arg0, arg1` | Script |
| `TempBGM` | `arg0` | Script |
| `ScrCmd_088` | `arg0` | Script |
| `ChatotHasCry` | `arg0` | Script |
| `ChatotStartRecording` | `arg0` | Script |
| `ChatotStopRecording` | `` | Script |
| `ChatotSaveRecording` | `` | Script |
| `ScrCmd_093` | `` | Script |
| `ApplyMovement` | `arg0, arg1` | Script |
| `WaitMovement` | `` | Script |
| `LockAll` | `` | Script |
| `ReleaseAll` | `` | Script |
| `Lock` | `arg0` | Script |
| `Release` | `arg0` | Script |
| `ShowPerson` | `arg0` | Script |
| `HidePerson` | `arg0` | Script |
| `ScrCmd_102` | `arg0, arg1` | Script |
| `ScrCmd_103` | `` | Script |
| `FacePlayer` | `` | Script |
| `GetPlayerCoords` | `arg0, arg1` | Script |
| `GetPersonCoords` | `arg0, arg1, arg2` | Script |
| `ScrCmd_107` | `arg0, arg1, arg2` | Script |
| `ScrCmd_108` | `arg0, arg1` | Script |
| `ScrCmd_109` | `arg0, arg1` | Script |
| `AddMoney` | `arg0` | Script |
| `SubMoneyImmediate` | `arg0` | Script |
| `HasEnoughMoneyImmediate` | `arg0, arg1` | Script |
| `ShowMoneyBox` | `arg0, arg1` | Script |
| `HideMoneyBox` | `` | Script |
| `UpdateMoneyBox` | `` | Script |
| `ScrCmd_116` | `arg0, arg1, arg2` | Script |
| `ScrCmd_117` | `` | Script |
| `ScrCmd_118` | `arg0` | Script |
| `GetCoinAmount` | `arg0` | Script |
| `GiveCoins` | `arg0` | Script |
| `TakeCoins` | `arg0` | Script |
| `GiveAthletePoints` | `arg0` | Script |
| `TakeAthletePoints` | `arg0` | Script |
| `CheckAthletePoints` | `arg0, arg1` | Script |
| `GiveItem` | `arg0, arg1, arg2` | Script |
| `TakeItem` | `arg0, arg1, arg2` | Script |
| `HasSpaceForItem` | `arg0, arg1, arg2` | Script |
| `HasItem` | `arg0, arg1, arg2` | Script |
| `ItemIsTMOrHM` | `arg0, arg1` | Script |
| `GetItemPocket` | `arg0, arg1` | Script |
| `SetStarterChoice` | `arg0` | Script |
| `GenderMsgBox` | `arg0, arg1` | Script |
| `GetSealQuantity` | `arg0, arg1` | Script |
| `GiveOrTakeSeal` | `arg0, arg1` | Script |
| `GiveRandomSeal` | `arg0, arg1, arg2` | Script |
| `ScrCmd_136` | `arg0, arg1` | Script |
| `GiveMon` | `species, level, heldItem, form, ability, retvar` | Script |
| `GiveEgg` | `arg0, arg1` | Script |
| `SetMonMove` | `arg0, arg1, arg2` | Script |
| `MonHasMove` | `arg0, arg1, arg2` | Script |
| `GetPartySlotWithMove` | `arg0, arg1` | Script |
| `GetPhoneBookRematch` | `arg0, arg1` | Script |
| `NameRival` | `arg0` | Script |
| `GetFriendSprite` | `arg0` | Script |
| `RegisterPokegearCard` | `arg0` | Script |
| `RegisterGearNumber` | `arg0` | Script |
| `CheckRegisteredPhoneNumber` | `arg0, arg1` | Script |
| `ScrCmd_148` | `arg0, arg1` | Script |
| `UnsetPhoneCallTrigger` | `arg0` | Script |
| `RestoreOverworld` | `` | Script |
| `ScrCmd_151` | `` | Script |
| `ScrCmd_152` | `` | Script |
| `ScrCmd_153` | `arg0` | Script |
| `ScrCmd_154` | `arg0, arg1, arg2` | Script |
| `ScrCmd_155` | `arg0, arg1` | Script |
| `ScrCmd_156` | `` | Script |
| `TownMap` | `` | Script |
| `ScrCmd_158` | `arg0` | Script |
| `ScrCmd_159` | `` | Script |
| `ScrCmd_160` | `` | Script |
| `ScrCmd_161` | `` | Script |
| `ScrCmd_162` | `` | Script |
| `HOFCredits` | `arg0` | Script |
| `ScrCmd_164` | `` | Script |
| `ScrCmd_165` | `arg0, arg1` | Script |
| `ScrCmd_166` | `arg0` | Script |
| `ChooseStarter` | `` | Script |
| `GetTrainerPathToPlayer` | `arg0` | Script |
| `TrainerStepTowardsPlayer` | `arg0, arg1` | Script |
| `GetTrainerEyeType` | `arg0` | Script |
| `GetEyeTrainerNum` | `arg0, arg1` | Script |
| `NamePlayer` | `arg0` | Script |
| `NicknameInput` | `arg0, arg1` | Script |
| `FadeScreen` | `arg0, speed, direction, color` | Script |
| `WaitFade` | `` | Script |
| `Warp` | `arg0, arg1=0, arg2, arg3, arg4` | Script |
| `RockClimb` | `arg0` | Script |
| `Surf` | `arg0` | Script |
| `Waterfall` | `arg0` | Script |
| `ScrCmd_180` | `arg0, arg1, arg2` | Script |
| `FlashEffect` | `` | Script |
| `Whirlpool` | `arg0` | Script |
| `ScrCmd_183` | `arg0` | Script |
| `PlayerOnBikeCheck` | `arg0` | Script |
| `PlayerOnBikeSet` | `arg0` | Script |
| `SetBikeStateLock` | `lock` | Script |
| `GetPlayerState` | `var` | Script |
| `SetAvatarBits` | `arg0` | Script |
| `UpdateAvatarState` | `` | Script |
| `BufferPlayersName` | `slot` | Script |
| `BufferRivalsName` | `slot` | Script |
| `BufferFriendsName` | `slot` | Script |
| `BufferMonSpeciesName` | `slot, party_pos` | Script |
| `BufferItemName` | `slot, item` | Script |
| `BufferPocketName` | `slot, pocket` | Script |
| `BufferTMHMMoveName` | `slot, tmhm` | Script |
| `BufferMoveName` | `slot, move` | Script |
| `BufferInt` | `slot, value` | Script |
| `BufferPartyMonNick` | `slot, party_pos` | Script |
| `BufferTrainerClassName` | `slot, trcls` | Script |
| `BufferPlayerUnionAvatarClassName` | `slot` | Script |
| `BufferSpeciesName` | `slot, species, arg2, arg3` | Script |
| `BufferStarterSpeciesName` | `slot` | Script |
| `BufferDPPtRivalStarterSpeciesName` | `slot` | Script |
| `BufferDPPtFriendStarterSpeciesName` | `slot` | Script |
| `GetStarterChoice` | `var` | Script |
| `BufferDecorationName` | `slot, deco` | Script |
| `ScrCmd_208` | `slot, unk` | Script |
| `ScrCmd_209` | `slot, unk` | Script |
| `BufferMapSecName` | `slot, location` | Script |
| `ScrCmd_211` | `arg0, arg1` | Script |
| `GetTrainerNum` | `arg0` | Script |
| `TrainerBattle` | `trainer, arg1, arg2, arg3` | Script |
| `TrainerMessage` | `trainer, param` | Script |
| `GetTrainerMsgParams` | `intro, after, _1poke` | Script |
| `GetRematchMsgParams` | `arg0, arg1, arg2` | Script |
| `TrainerIsDoubleBattle` | `arg0` | Script |
| `EncounterMusic` | `arg0` | Script |
| `WhiteOut` | `` | Script |
| `CheckBattleWon` | `var` | Script |
| `StaticWildWonOrCaughtCheck` | `arg0, arg1` | Script |
| `PartyCheckForDouble` | `arg0` | Script |
| `ScrCmd_223` | `` | Script |
| `ScrCmd_224` | `` | Script |
| `GoToIfTrainerDefeated` | `arg0` | Script |
| `ScrCmd_226` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_227` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_228` | `arg0` | Script |
| `ScrCmd_229` | `arg0` | Script |
| `ScrCmd_230` | `` | Script |
| `ScrCmd_231` | `` | Script |
| `ScrCmd_232` | `arg0` | Script |
| `ScrCmd_233` | `arg0` | Script |
| `ScrCmd_234` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_235` | `arg0` | Script |
| `ScrCmd_236` | `arg0` | Script |
| `ScrCmd_237` | `` | Script |
| `PartyHasPokerus` | `arg0` | Script |
| `MonGetGender` | `arg0, arg1` | Script |
| `SetDynamicWarp` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `GetDynamicWarpFloorNo` | `arg0` | Script |
| `ElevatorCurFloorBox` | `arg0, arg1, arg2, arg3` | Script |
| `CountJohtoDexSeen` | `var` | Script |
| `CountJohtoDexOwned` | `var` | Script |
| `CountNationalDexSeen` | `var` | Script |
| `CountNationalDexOwned` | `var` | Script |
| `ScrCmd_247` | `` | Script |
| `GetDexEvalResult` | `national, msg, fanfare` | Script |
| `RocketTrapBattle` | `arg0, arg1` | Script |
| `ScrCmd_250` | `arg0, arg1` | Script |
| `CatchingTutorial` | `` | Script |
| `ScrCmd_252` | `` | Script |
| `GetSaveFileState` | `arg0` | Script |
| `SaveGameNormal` | `arg0` | Script |
| `ScrCmd_255` | `arg0, arg1` | Script |
| `ScrCmd_256` | `arg0` | Script |
| `ScrCmd_257` | `arg0` | Script |
| `ScrCmd_258` | `` | Script |
| `ScrCmd_259` | `arg0` | Script |
| `ScrCmd_260` | `arg0` | Script |
| `ScrCmd_261` | `arg0` | Script |
| `ScrCmd_262` | `` | Script |
| `ScrCmd_263` | `` | Script |
| `ScrCmd_264` | `arg0` | Script |
| `ScrCmd_265` | `` | Script |
| `ScrCmd_266` | `` | Script |
| `ScrCmd_267` | `arg0, arg1` | Script |
| `ScrCmd_268` | `arg0` | Script |
| `ScrCmd_269` | `arg0` | Script |
| `ScrCmd_270` | `` | Script |
| `ScrCmd_271` | `arg0, arg1` | Script |
| `ScrCmd_272` | `arg0` | Script |
| `ScrCmd_273` | `arg0` | Script |
| `ScrCmd_274` | `arg0, arg1` | Script |
| `MartBuy` | `arg0` | Script |
| `SpecialMartBuy` | `arg0` | Script |
| `DecorationMart` | `arg0` | Script |
| `SealMart` | `arg0` | Script |
| `OverworldWhiteOut` | `` | Script |
| `SetSpawn` | `arg0` | Script |
| `GetPlayerGender` | `var` | Script |
| `HealParty` | `` | Script |
| `ScrCmd_283` | `` | Script |
| `ScrCmd_284` | `` | Script |
| `ScrCmd_285` | `arg0` | Script |
| `ScrCmd_286` | `` | Script |
| `BufferUnionRoomAvatarChoices` | `` | Script |
| `UnionRoomAvatarIdxToTrainerClass` | `arg0, arg1` | Script |
| `ScrCmd_289` | `arg0` | Script |
| `CheckPokedex` | `arg0` | Script |
| `GivePokedex` | `` | Script |
| `CheckRunningShoes` | `arg0` | Script |
| `GiveRunningShoes` | `` | Script |
| `CheckBadge` | `badge, var` | Script |
| `GiveBadge` | `badge` | Script |
| `CountBadges` | `var` | Script |
| `ScrCmd_297` | `arg0` | Script |
| `ScrCmd_298` | `` | Script |
| `CheckEscortMode` | `arg0` | Script |
| `SetEscortMode` | `` | Script |
| `ClearEscortMode` | `` | Script |
| `CheckStepTakenFlag` | `arg0` | Script |
| `SetStepTakenFlag` | `` | Script |
| `GetStepTakenFlag` | `` | Script |
| `CheckGameClearFlag` | `arg0` | Script |
| `SetGameClearFlag` | `` | Script |
| `ScrCmd_307` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_308` | `arg0` | Script |
| `ScrCmd_309` | `arg0` | Script |
| `ScrCmd_310` | `arg0` | Script |
| `ScrCmd_311` | `arg0` | Script |
| `BufferDaycareMonNicks` | `` | Script |
| `GetDaycareState` | `arg0` | Script |
| `EcruteakGymInit` | `` | Script |
| `BindEcruteakGymCandleToTrainerSprite` | `` | Script |
| `UnbindEcruteakGymCandleFromTrainerSprite` | `` | Script |
| `ExtinguishEcruteakGymCandle` | `arg0` | Script |
| `CianwoodGymInit` | `` | Script |
| `CianwoodGymTurnWinch` | `arg0` | Script |
| `VermilionGymInit` | `` | Script |
| `VermilionGymLockAction` | `arg0, arg1` | Script |
| `VermilionGymCanCheck` | `arg0, arg1` | Script |
| `ResampleVermilionGymCans` | `` | Script |
| `VioletGymInit` | `` | Script |
| `VioletGymElevator` | `` | Script |
| `AzaleaGymInit` | `` | Script |
| `AzaleaGymSpinarak` | `arg0` | Script |
| `AzaleaGymSwitch` | `arg0` | Script |
| `BlackthornGymInit` | `` | Script |
| `FuchsiaGymInit` | `` | Script |
| `ViridianGymInit` | `` | Script |
| `GetPartyCount` | `var` | Script |
| `ScrCmd_333` | `arg0` | Script |
| `ScrCmd_334` | `arg0` | Script |
| `ScrCmd_335` | `arg0, arg1` | Script |
| `BufferBerryName` | `slot, item, quantity` | Script |
| `BufferNatureName` | `slot, nature` | Script |
| `MovePerson` | `arg0, arg1, arg2` | Script |
| `MovePersonFacing` | `person, x, z, y, facing` | Script |
| `SetObjectMovementType` | `arg0, arg1` | Script |
| `SetObjectFacing` | `arg0, arg1` | Script |
| `MoveWarp` | `arg0, arg1, arg2` | Script |
| `MoveBGEvent` | `arg0, arg1, arg2` | Script |
| `ScrCmd_344` | `arg0, arg1` | Script |
| `AddWaitingIcon` | `` | Script |
| `RemoveWaitingIcon` | `` | Script |
| `ScrCmd_347` | `arg0` | Script |
| `WaitButtonOrDelay` | `arg0` | Script |
| `PartySelectUI` | `` | Script |
| `ScrCmd_350` | `` | Script |
| `GetPartySelection` | `arg0` | Script |
| `PokemonSummaryScreen` | `arg0, arg1, arg2` | Script |
| `GetMoveSelection` | `arg0, arg1` | Script |
| `GetPartyMonSpecies` | `slot, var` | Script |
| `PartyMonIsMine` | `slot, var` | Script |
| `PartyCountNotEgg` | `var` | Script |
| `CountAliveMons` | `var, except=PARTY_SIZE` | Script |
| `CountAliveMonsAndPC` | `arg0` | Script |
| `PartyCountEgg` | `arg0` | Script |
| `SubMoneyVar` | `var` | Script |
| `RetrieveDaycareMon` | `arg0, arg1` | Script |
| `GiveLoanMon` | `arg0, arg1, arg2` | Script |
| `CheckReturnLoanMon` | `arg0, arg1, arg2` | Script |
| `ReturnLoanMon` | `arg0` | Script |
| `ResetDaycareEgg` | `` | Script |
| `GiveDaycareEgg` | `` | Script |
| `BufferDaycareWithdrawCost` | `arg0, arg1` | Script |
| `HasEnoughMoneyVar` | `var, amount` | Script |
| `EggHatchAnim` | `` | Script |
| `ScrCmd_370` | `arg0` | Script |
| `BufferDaycareMonGrowth` | `arg0, arg1` | Script |
| `GetTailDaycareMonSpeciesAndNick` | `arg0` | Script |
| `PutMonInDaycare` | `arg0` | Script |
| `ScrCmd_374` | `arg0` | Script |
| `MakeObjectVisible` | `arg0` | Script |
| `ScrCmd_376` | `` | Script |
| `ScrCmd_377` | `arg0` | Script |
| `ViewRankings` | `scope, page, record` | Script |
| `ScrCmd_379` | `arg0` | Script |
| `Random` | `arg0, arg1` | Script |
| `ScrCmd_381` | `arg0, arg1` | Script |
| `MonGetFriendship` | `arg0, arg1` | Script |
| `MonAddFriendship` | `arg0, arg1` | Script |
| `MonSubtractFriendship` | `arg0, arg1` | Script |
| `BufferDaycareMonStats` | `arg0, arg1, arg2, arg3` | Script |
| `GetPlayerFacing` | `arg0` | Script |
| `GetDaycareCompatibility` | `arg0` | Script |
| `CheckDaycareEgg` | `arg0` | Script |
| `PlayerHasSpecies` | `arg0, arg1` | Script |
| `SizeRecordCompare` | `arg0, arg1` | Script |
| `SizeRecordUpdate` | `arg0` | Script |
| `BufferMonSize` | `arg0, arg1, arg2` | Script |
| `BufferRecordSize` | `arg0, arg1, arg2` | Script |
| `ScrCmd_394` | `arg0` | Script |
| `ScrCmd_395` | `arg0` | Script |
| `CountMonMoves` | `arg0, arg1` | Script |
| `MonForgetMove` | `arg0, arg1` | Script |
| `MonGetMove` | `arg0, arg1, arg2` | Script |
| `BufferPartyMonMoveName` | `slot, party_pos, move_pos` | Script |
| `StrengthFlagAction` | `action, var=0` | Script |
| `FlashAction` | `action, var=0` | Script |
| `DefogAction` | `action, var=0` | Script |
| `ScrCmd_403` | `arg0, arg1` | Script |
| `ScrCmd_404` | `arg0, arg1, arg2` | Script |
| `ScrCmd_405` | `arg0, arg1, arg2` | Script |
| `ScrCmd_406` | `arg0` | Script |
| `ScrCmd_407` | `arg0, arg1` | Script |
| `ScrCmd_408` | `arg0, arg1` | Script |
| `ScrCmd_409` | `` | Script |
| `ScrCmd_410` | `arg0, arg1` | Script |
| `ScrCmd_411` | `` | Script |
| `ScrCmd_412` | `arg0, arg1, arg2` | Script |
| `ScrCmd_413` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_414` | `arg0` | Script |
| `ScrCmd_415` | `arg0` | Script |
| `ScrCmd_416` | `arg0, arg1, arg2` | Script |
| `ScrCmd_417` | `arg0, arg1` | Script |
| `ScrCmd_418` | `arg0, arg1` | Script |
| `ScrCmd_419` | `arg0` | Script |
| `ScrCmd_420` | `arg0` | Script |
| `ScrCmd_421` | `arg0, arg1, arg2` | Script |
| `ScrCmd_422` | `arg0, arg1, arg2, arg3` | Script |
| `CheckJohtoDexComplete` | `arg0` | Script |
| `CheckNationalDexComplete` | `arg0` | Script |
| `ShowCertificate` | `arg0` | Script |
| `KenyaCheck` | `arg0, arg1, arg2` | Script |
| `ScrCmd_427` | `arg0` | Script |
| `MonGiveMail` | `arg0` | Script |
| `CountFossils` | `arg0` | Script |
| `SetPhoneCall` | `arg0, arg1, arg2` | Script |
| `RunPhoneCall` | `` | Script |
| `GetFossilPokemon` | `arg0, arg1` | Script |
| `GetFossilMinimumAmount` | `arg0, arg1, arg2` | Script |
| `PartyCountMonsAtOrBelowLevel` | `arg0, arg1` | Script |
| `SurvivePoisoning` | `arg0, arg1` | Script |
| `ScrCmd_436` | `` | Script |
| `DebugWatch` | `arg0` | Script |
| `GetStdMsgNaix` | `arg0, arg1` | Script |
| `NonNPCMsgExtern` | `arg0, arg1` | Script |
| `MsgBoxExtern` | `arg0, arg1` | Script |
| `ScrCmd_441` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_442` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_443` | `arg0` | Script |
| `ScrCmd_444` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_445` | `arg0` | Script |
| `ScrCmd_446` | `arg0` | Script |
| `SafariZoneAction` | `arg0, arg1` | Script |
| `ScrCmd_448` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_449` | `` | Script |
| `ScrCmd_450` | `` | Script |
| `ScrCmd_451` | `arg0` | Script |
| `ScrCmd_452` | `arg0, arg1` | Script |
| `ScrCmd_453` | `` | Script |
| `ScrCmd_454` | `` | Script |
| `ScrCmd_455` | `` | Script |
| `ScrCmd_456` | `arg0` | Script |
| `MonGetNature` | `arg0, arg1` | Script |
| `GetPartySlotWithNature` | `arg0, arg1` | Script |
| `ScrCmd_459` | `` | Script |
| `LoadPhoneDat` | `arg0, arg1` | Script |
| `GetPhoneContactMsgIds` | `arg0, arg1, arg2` | Script |
| `ScrCmd_462` | `arg0` | Script |
| `EnableMassOutbreaks` | `` | Script |
| `CreateRoamer` | `arg0` | Script |
| `ScrCmd_465` | `arg0, arg1=0, arg2=0` | Script |
| `ScrCmd_466` | `arg0, arg1` | Script |
| `MoveRelearnerInit` | `arg0` | Script |
| `MoveTutorInit` | `arg0, arg1` | Script |
| `MoveRelearnerGetResult` | `arg0` | Script |
| `LoadNPCTrade` | `arg0` | Script |
| `GetOfferedSpecies` | `arg0` | Script |
| `NPCTradeGetReqSpecies` | `arg0` | Script |
| `NPCTradeExec` | `arg0` | Script |
| `NPCTradeEnd` | `` | Script |
| `ScrCmd_475` | `` | Script |
| `EnablePokedexFormDetection` | `` | Script |
| `NatDexFlagAction` | `arg0, arg1` | Script |
| `MonGetRibbonCount` | `arg0, arg1` | Script |
| `GetPartyRibbonCount` | `arg0` | Script |
| `MonHasRibbon` | `arg0, arg1, arg2` | Script |
| `GiveRibbon` | `arg0, arg1` | Script |
| `BufferRibbonName` | `arg0, arg1` | Script |
| `GetEVTotal` | `arg0, arg1` | Script |
| `GetWeekday` | `arg0` | Script |
| `StartBattleRegulationMenuTask` | `arg0` | Script |
| `Dummy486` | `` | Script |
| `PokeCenAnim` | `arg0` | Script |
| `ElevatorAnim` | `direction, distance` | Script |
| `MysteryGift` | `arg0, arg1=0, arg2=0` | Script |
| `NopVar490` | `arg0` | Script |
| `ScrCmd_491` | `arg0` | Script |
| `ScrCmd_492` | `arg0, arg1, arg2` | Script |
| `PromptEasyChat` | `arg0, arg1, arg2` | Script |
| `ScrCmd_494` | `arg0, arg1` | Script |
| `GetGameVersion` | `arg0` | Script |
| `GetPartyLead` | `arg0` | Script |
| `GetMonTypes` | `arg0, arg1, arg2` | Script |
| `PrimoPasswordCheck1` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `PrimoPasswordCheck2` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_500` | `arg0` | Script |
| `ScrCmd_501` | `arg0` | Script |
| `ScrCmd_502` | `arg0` | Script |
| `LotoIDGet` | `arg0` | Script |
| `LotoIDSearch` | `arg0, arg1, arg2, arg3` | Script |
| `LotoIDSet` | `` | Script |
| `BufferBoxMonNick` | `arg0, arg1` | Script |
| `CountPCEmptySpace` | `arg0` | Script |
| `PalParkAction` | `arg0` | Script |
| `ScrCmd_509` | `arg0` | Script |
| `ScrCmd_510` | `` | Script |
| `PalParkScoreGet` | `arg0, arg1` | Script |
| `PlayerMovementSavingSet` | `` | Script |
| `PlayerMovementSavingClear` | `` | Script |
| `HallOfFameAnim` | `num` | Script |
| `AddSpecialGameStat` | `arg0` | Script |
| `BufferFashionName` | `arg0, arg1` | Script |
| `ScrCmd_517` | `arg0, arg1` | Script |
| `ScrCmd_518` | `arg0` | Script |
| `ScrCmd_519` | `arg0` | Script |
| `ScrCmd_520` | `` | Script |
| `ScrCmd_521` | `` | Script |
| `ScrCmd_522` | `arg0` | Script |
| `ScrCmd_523` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_524` | `arg0, arg1, arg2` | Script |
| `ScrCmd_525` | `arg0` | Script |
| `ScrCmd_526` | `arg0` | Script |
| `ScrCmd_527` | `arg0` | Script |
| `ScrCmd_528` | `arg0` | Script |
| `GetPartyLeadAlive` | `arg0` | Script |
| `ScrCmd_530` | `arg0, arg1` | Script |
| `BufferBackgroundName` | `arg0, arg1` | Script |
| `CheckCoinsImmediate` | `arg0, arg1` | Script |
| `CheckGiveCoins` | `arg0, arg1` | Script |
| `ScrCmd_534` | `arg0` | Script |
| `MonGetLevel` | `arg0, arg1` | Script |
| `ScrCmd_536` | `arg0, arg1` | Script |
| `ScrCmd_537` | `` | Script |
| `ScrCmd_538` | `arg0, arg1` | Script |
| `ScrCmd_539` | `arg0` | Script |
| `ScrCmd_540` | `arg0` | Script |
| `BufferIntEx` | `arg0, arg1, arg2, arg3` | Script |
| `MonGetContestValue` | `arg0, arg1, arg2` | Script |
| `ScrCmd_543` | `arg0` | Script |
| `ScrCmd_544` | `arg0, arg1` | Script |
| `ScrCmd_545` | `arg0` | Script |
| `ScrCmd_546` | `arg0, arg1` | Script |
| `ScrCmd_547` | `arg0` | Script |
| `ScrCmd_548` | `` | Script |
| `ScrCmd_549` | `arg0` | Script |
| `ScrCmd_550` | `arg0` | Script |
| `ScrCmd_551` | `arg0` | Script |
| `ScrCmd_552` | `arg0, arg1` | Script |
| `ScrCmd_553` | `arg0, arg1` | Script |
| `ScrCmd_554` | `arg0` | Script |
| `ScrCmd_555` | `arg0` | Script |
| `ScrCmd_556` | `arg0` | Script |
| `CheckBattlePoints` | `arg0, arg1` | Script |
| `UnionRoomAvatarIdxToSprite` | `arg0, arg1` | Script |
| `ScrCmd_559` | `arg0, arg1` | Script |
| `ScrCmd_560` | `arg0, arg1` | Script |
| `ScreenShake` | `arg0, arg1, arg2, arg3` | Script |
| `MultiBattle` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_563` | `arg0, arg1, arg2` | Script |
| `ScrCmd_564` | `arg0` | Script |
| `ScrCmd_565` | `arg0` | Script |
| `ScrCmd_566` | `` | Script |
| `GetDPPlPrizeItemIDAndCost` | `arg0, arg1, arg2` | Script |
| `ScrCmd_568` | `arg0, arg1` | Script |
| `ScrCmd_569` | `arg0` | Script |
| `CheckCoinsVar` | `arg0, arg1` | Script |
| `ScrCmd_571` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `GetUniqueSealsQuantity` | `arg0` | Script |
| `ScrCmd_573` | `` | Script |
| `ScrCmd_574` | `arg0, arg1` | Script |
| `ScrCmd_575` | `arg0, arg1` | Script |
| `ScrCmd_576` | `arg0` | Script |
| `ScrCmd_577` | `` | Script |
| `ScrCmd_578` | `` | Script |
| `ScrCmd_579` | `` | Script |
| `BufferSealName` | `arg0, arg1` | Script |
| `LockLastTalked` | `` | Script |
| `ScrCmd_582` | `arg0, arg1, arg2` | Script |
| `ScrCmd_583` | `arg0, arg1` | Script |
| `PartyLegalCheck` | `arg0` | Script |
| `ScrCmd_585` | `` | Script |
| `ScrCmd_586` | `arg0` | Script |
| `ScrCmd_587` | `` | Script |
| `LatiCaughtCheck` | `arg0` | Script |
| `WildBattle` | `species, level, shiny=0` | Script |
| `GetTrcardStars` | `arg0` | Script |
| `ScrCmd_591` | `` | Script |
| `ScrCmd_592` | `arg0` | Script |
| `ShowSaveStats` | `` | Script |
| `HideSaveStats` | `` | Script |
| `ScrCmd_595` | `arg0` | Script |
| `ScrCmd_596` | `arg0` | Script |
| `ScrCmd_597` | `` | Script |
| `ScrCmd_598` | `arg0` | Script |
| `ScrCmd_599` | `` | Script |
| `ScrCmd_600` | `` | Script |
| `FollowMonFacePlayer` | `` | Script |
| `ToggleFollowingPokemonMovement` | `flag` | Script |
| `WaitFollowingPokemonMovement` | `` | Script |
| `FollowingPokemonMovement` | `movement` | Script |
| `ScrCmd_605` | `arg0, arg1` | Script |
| `ScrCmd_606` | `` | Script |
| `ScrCmd_607` | `` | Script |
| `ScrCmd_608` | `` | Script |
| `ScrCmd_609` | `` | Script |
| `ScrCmd_610` | `arg0` | Script |
| `Pokeathlon` | `arg0, arg1, arg2, arg3, arg4, arg5, arg6` | Script |
| `GetNPCTradeUnusedFlag` | `arg0` | Script |
| `GetPhoneContactRandomGiftBerry` | `arg0` | Script |
| `GetPhoneContactGiftItem` | `arg0` | Script |
| `CameronPhoto` | `arg0` | Script |
| `CountSavedPhotos` | `arg0` | Script |
| `OpenPhotoAlbum` | `` | Script |
| `PhotoAlbumIsFull` | `arg0` | Script |
| `RocketCostumeFlagCheck` | `arg0` | Script |
| `RocketCostumeFlagAction` | `arg0` | Script |
| `PlaceStarterBallsInElmsLab` | `` | Script |
| `ScrCmd_622` | `arg0, arg1` | Script |
| `AnimApricornTree` | `arg0` | Script |
| `ApricornTreeGetApricorn` | `arg0` | Script |
| `GiveApricornFromTree` | `arg0, arg1, arg2` | Script |
| `BufferApricornName` | `arg0, arg1` | Script |
| `ScrCmd_627` | `arg0` | Script |
| `ScrCmd_628` | `arg0, arg1` | Script |
| `ScrCmd_629` | `` | Script |
| `ScrCmd_630` | `arg0` | Script |
| `ScrCmd_631` | `arg0, arg1, arg2` | Script |
| `CountPartyMonsOfSpecies` | `arg0, arg1` | Script |
| `ScrCmd_633` | `arg0, arg1, arg2` | Script |
| `ScrCmd_634` | `arg0, arg1` | Script |
| `ScrCmd_635` | `arg0, arg1` | Script |
| `ScrCmd_636` | `arg0` | Script |
| `ScrCmd_637` | `arg0, arg1, arg2` | Script |
| `ScrCmd_638` | `arg0, arg1, arg2` | Script |
| `ScrCmd_639` | `arg0, arg1, arg2` | Script |
| `ScrCmd_640` | `arg0` | Script |
| `SaveWipeExtraChunks` | `` | Script |
| `ScrCmd_642` | `arg0` | Script |
| `ScrCmd_643` | `arg0, arg1, arg2` | Script |
| `ScrCmd_644` | `arg0, arg1, arg2` | Script |
| `ScrCmd_645` | `arg0, arg1, arg2` | Script |
| `ScrCmd_646` | `arg0` | Script |
| `GetPartySlotWithSpecies` | `arg0, arg1` | Script |
| `ScrCmd_648` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScratchOffCard` | `` | Script |
| `ScratchOffCardEnd` | `` | Script |
| `GetScratchOffPrize` | `arg0, arg1, arg2` | Script |
| `ScrCmd_652` | `arg0, arg1, arg2` | Script |
| `MoveTutorChooseMove` | `mon_slot, tutor_no, arg2, arg3` | Script |
| `TutorMoveTeachInSlot` | `party_slot, move_slot, move` | Script |
| `TutorMoveGetPrice` | `arg0, arg1` | Script |
| `ScrCmd_656` | `arg0, arg1` | Script |
| `StatJudge` | `mon_slot, iv_sum, best_stat, best_iv` | Script |
| `BufferStatName` | `arg0, arg1` | Script |
| `SetMonForm` | `arg0, arg1` | Script |
| `BufferTrainerName` | `arg0, arg1` | Script |
| `ScrCmd_661` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_662` | `arg0, arg1, arg2` | Script |
| `ScrCmd_663` | `arg0` | Script |
| `ScrCmd_664` | `` | Script |
| `ScrCmd_665` | `arg0` | Script |
| `ScrCmd_666` | `arg0` | Script |
| `ScrCmd_667` | `arg0` | Script |
| `BufferTypeName` | `arg0, arg1` | Script |
| `GetItemQuantity` | `arg0, arg1` | Script |
| `GetHiddenPowerType` | `arg0, arg1` | Script |
| `SetFavoriteMon` | `` | Script |
| `GetFavoriteMon` | `arg0, arg1, arg2` | Script |
| `GetOwnedRotomForms` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `CountTranformedRotomsInParty` | `arg0, arg1` | Script |
| `UpdateRotomForm` | `arg0, arg1, arg2, arg3` | Script |
| `GetPartyMonForm` | `arg0, arg1` | Script |
| `ScrCmd_677` | `arg0, arg1` | Script |
| `ScrCmd_678` | `arg0, arg1` | Script |
| `ScrCmd_679` | `` | Script |
| `AddSpecialGameStat2` | `arg0` | Script |
| `ScrCmd_681` | `arg0` | Script |
| `ScrCmd_682` | `arg0` | Script |
| `GetStaticEncounterOutcome` | `arg0` | Script |
| `ScrCmd_684` | `arg0` | Script |
| `GetPlayerXYZ` | `arg0, arg1, arg2` | Script |
| `ScrCmd_686` | `arg0, arg1` | Script |
| `ScrCmd_687` | `arg0` | Script |
| `GetPartySlotWithFatefulEncounter` | `arg0, arg1` | Script |
| `CommSanitizeParty` | `var_result` | Script |
| `DaycareSanitizeMon` | `arg0, arg1` | Script |
| `ScrCmd_691` | `arg0` | Script |
| `BufferBattleHallStreak` | `arg0, arg1, arg2, arg3, arg4, arg5` | Script |
| `BattleHallCountUsedSpecies` | `arg0` | Script |
| `BattleHallGetTotalStreak` | `arg0` | Script |
| `ScrCmd_695` | `arg0` | Script |
| `ScrCmd_696` | `arg0` | Script |
| `ScrCmd_697` | `arg0` | Script |
| `FollowerPokeIsEventTrigger` | `arg0, arg1, arg2` | Script |
| `ScrCmd_699` | `` | Script |
| `ScrCmd_700` | `` | Script |
| `MonHasItem` | `arg0, arg1` | Script |
| `BattleTowerSetUpMultiBattle` | `` | Script |
| `SetPlayerVolume` | `arg0` | Script |
| `ScrCmd_704` | `arg0, arg1` | Script |
| `ScrCmd_705` | `arg0, arg1` | Script |
| `ScrCmd_706` | `arg0` | Script |
| `CheckMonSeen` | `arg0, arg1` | Script |
| `ScrCmd_708` | `arg0` | Script |
| `ScrCmd_709` | `` | Script |
| `ScrCmd_710` | `` | Script |
| `FollowMonInteract` | `` | Script |
| `ScrCmd_712` | `arg0` | Script |
| `AlphPuzzle` | `arg0` | Script |
| `OpenAlphHiddenRoom` | `arg0` | Script |
| `UpdateDaycareMonObjects` | `` | Script |
| `ScrCmd_716` | `` | Script |
| `ScrCmd_717` | `arg0` | Script |
| `ScrCmd_718` | `arg0, arg1` | Script |
| `ScrCmd_719` | `arg0, arg1` | Script |
| `ScrCmd_720` | `arg0` | Script |
| `ScrCmd_721` | `arg0` | Script |
| `ScrCmd_722` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_723` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `ScrCmd_724` | `arg0, arg1` | Script |
| `ScrCmd_725` | `arg0, arg1` | Script |
| `ProcessSoundplate` | `` | Script |
| `GetFollowPokePartyIndex` | `arg0` | Script |
| `ScrCmd_728` | `arg0, arg1` | Script |
| `ScrCmd_729` | `arg0` | Script |
| `ScrCmd_730` | `arg0` | Script |
| `ScrCmd_731` | `` | Script |
| `ScrCmd_732` | `arg0` | Script |
| `ScrCmd_733` | `arg0, arg1` | Script |
| `ScrCmd_734` | `arg0` | Script |
| `ScrCmd_735` | `arg0` | Script |
| `ClearKurtApricorn` | `` | Script |
| `ScrCmd_737` | `arg0` | Script |
| `GetTotalApricornCount` | `arg0` | Script |
| `ScrCmd_739` | `` | Script |
| `ScrCmd_740` | `arg0, arg1` | Script |
| `ScrCmd_741` | `arg0, arg1, arg2, arg3` | Script |
| `ScrCmd_742` | `arg0, arg1, arg2` | Script |
| `ScrCmd_743` | `arg0` | Script |
| `CreatePokeathlonFriendshipRoomStatues` | `` | Script |
| `BufferPokeathlonCourseName` | `arg0, arg1` | Script |
| `TouchscreenMenuHide` | `` | Script |
| `TouchscreenMenuShow` | `` | Script |
| `GetMenuChoice` | `arg0` | Script |
| `MenuInitStdGmm` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `MenuInit` | `arg0, arg1, arg2, arg3, arg4` | Script |
| `MenuItemAdd` | `arg0, arg1, arg2` | Script |
| `MenuExec` | `` | Script |
| `RockSmashItemCheck` | `arg0, arg1, arg2` | Script |
| `TryHeadbuttEncounter` | `arg0` | Script |
| `LegendCutsceneClearBellAnimBegin` | `` | Script |
| `LegendCutsceneClearBellAnimEnd` | `` | Script |
| `LegendCutsceneClearBellRiseFromBag` | `` | Script |
| `LegendCutsceneClearBellShimmer` | `arg0` | Script |
| `LegendCutsceneLugiaEyeGlimmerEffect` | `` | Script |
| `ScrCmd_760` | `` | Script |
| `LegendCutsceneMoveCameraTo` | `arg0` | Script |
| `LegendCutscenePanCameraTo` | `arg0` | Script |
| `LegendCutsceneWaitCameraPan` | `` | Script |
| `LegendCutsceneBirdFinalApproach` | `` | Script |
| `LegendCutsceneWavesOrLeavesEffectBegin` | `` | Script |
| `LegendCutsceneWavesOrLeavesEffectEnd` | `` | Script |
| `LegendCutsceneLugiaArrivesEffectBegin` | `` | Script |
| `LegendCutsceneLugiaArrivesEffectEnd` | `` | Script |
| `LegendCutsceneLugiaArrivesEffectCameraPan` | `` | Script |
| `CheckSeenAllLetterUnown` | `arg0` | Script |
| `ScrCmd_771` | `` | Script |
| `ScrCmd_772` | `` | Script |
| `Cinematic` | `arg0` | Script |
| `ShowLegendaryWing` | `arg0` | Script |
| `ScrCmd_775` | `arg0, arg1` | Script |
| `GiveTogepiEgg` | `` | Script |
| `ScrCmd_777` | `arg0, arg1` | Script |
| `GiveSpikyEarPichu` | `` | Script |
| `RadioMusicIsPlaying` | `arg0, arg1` | Script |
| `CasinoGame` | `arg0, arg1` | Script |
| `KenyaCheckPartyOrMailbox` | `arg0` | Script |
| `MartSell` | `` | Script |
| `SetFollowMonInhibitState` | `arg0` | Script |
| `ScriptOverlayCmd` | `arg0, arg1` | Script |
| `BugContestAction` | `arg0, arg1` | Script |
| `BufferBugContestWinner` | `arg0` | Script |
| `JudgeBugContest` | `arg0, arg1, arg2` | Script |
| `BufferBugContestMonNick` | `arg0, arg1` | Script |
| `BugContestGetTimeLeft` | `arg0` | Script |
| `IsBugContestantRegistered` | `arg0, arg1` | Script |
| `CheckSafariZoneChallengeCompleted` | `arg0, arg1` | Script |
| `UpdateSafariZoneIGT` | `` | Script |
| `BankTransaction` | `arg0, arg1` | Script |
| `CheckBankBalance` | `arg0, arg1` | Script |
| `ScrCmd_795` | `arg0, arg1` | Script |
| `ScrCmd_796` | `` | Script |
| `ScrCmd_797` | `` | Script |
| `BufferRulesetName` | `arg0` | Script |
| `ScrCmd_799` | `arg0` | Script |
| `ScrCmd_800` | `arg0` | Script |
| `ScrCmd_801` | `arg0` | Script |
| `ScrCmd_802` | `` | Script |
| `ScrCmd_803` | `arg0, arg1` | Script |
| `ScrCmd_804` | `arg0` | Script |
| `ScrCmd_805` | `` | Script |
| `ScrCmd_806` | `` | Script |
| `SetTrainerHouseSprite` | `arg0, arg1` | Script |
| `ScrCmd_808` | `arg0` | Script |
| `ShowTrainerHouseIntroMessage` | `arg0` | Script |
| `ScrCmd_810` | `` | Script |
| `ScrCmd_811` | `arg0, arg1` | Script |
| `ScrCmd_812` | `` | Script |
| `MomGiftCheck` | `arg0` | Script |
| `ScrCmd_814` | `` | Script |
| `ScrCmd_815` | `arg0` | Script |
| `UnownCircle` | `` | Script |
| `ScrCmd_817` | `arg0` | Script |
| `MystriStageGymmickInit` | `` | Script |
| `ScrCmd_819` | `` | Script |
| `ScrCmd_820` | `arg0` | Script |
| `GetBuenasPassword` | `arg0, arg1` | Script |
| `ScrCmd_822` | `` | Script |
| `ScrCmd_823` | `arg0` | Script |
| `ScrCmd_824` | `arg0` | Script |
| `GetShinyLeafCount` | `arg0, arg1` | Script |
| `TryGiveShinyLeafCrown` | `arg0` | Script |
| `GetPartyMonForm2` | `arg0, arg1` | Script |
| `MonAddContestValue` | `arg0, arg1, arg2` | Script |
| `ScrCmd_829` | `arg0` | Script |
| `ScrCmd_830` | `arg0` | Script |
| `ScrCmd_831` | `arg0` | Script |
| `ScrCmd_832` | `arg0` | Script |
| `ScrCmd_833` | `arg0` | Script |
| `ScrCmd_834` | `arg0` | Script |
| `ScrCmd_835` | `arg0` | Script |
| `CheckKyogreGroudonInParty` | `arg0` | Script |
| `ScrCmd_837` | `arg0` | Script |
| `BankOrWalletIsFull` | `arg0, arg1` | Script |
| `SysSetSleepFlag` | `arg0` | Script |
| `ScrCmd_840` | `arg0, arg1` | Script |
| `ScrCmd_841` | `arg0` | Script |
| `ScrCmd_842` | `arg0` | Script |
| `BufferItemNameIndef` | `arg0, arg1` | Script |
| `BufferItemNamePlural` | `arg0, arg1` | Script |
| `BufferPartyMonSpeciesNameIndef` | `arg0, arg1` | Script |
| `BufferSpeciesNameIndef` | `arg0, arg1, arg2, arg3` | Script |
| `BufferDPPtFriendStarterSpeciesNameIndef` | `arg0` | Script |
| `BufferFashionNameIndef` | `arg0, arg1` | Script |
| `BufferTrainerClassNameIndef` | `arg0, arg1` | Script |
| `BufferSealNamePlural` | `arg0, arg1` | Script |
| `Capitalize` | `arg0` | Script |
| `BufferDeptStoreFloorNo` | `arg0, arg1` | Script |
| `MapScript` | `kind, script` | Script |
| `MapScript2` | `var, val, script` | Script |
| `Compare` | `var, arg` | Script |
| `GoToIfUnset` | `flag, dest` | Script |
| `GoToIfSet` | `flag, dest` | Script |
| `GoToIfLt` | `dest` | Script |
| `GoToIfEq` | `dest` | Script |
| `GoToIfGt` | `dest` | Script |
| `GoToIfLe` | `dest` | Script |
| `GoToIfGe` | `dest` | Script |
| `GoToIfNe` | `dest` | Script |
| `CallIfUnset` | `flag, dest` | Script |
| `CallIfSet` | `flag, dest` | Script |
| `CallIfLt` | `dest` | Script |
| `CallIfEq` | `dest` | Script |
| `CallIfGt` | `dest` | Script |
| `CallIfLe` | `dest` | Script |
| `CallIfGe` | `dest` | Script |
| `CallIfNe` | `dest` | Script |
| `GoToIfDefeated` | `trainer, dest` | Script |
| `GoToIfNotDefeated` | `trainer, dest` | Script |
| `CallIfDefeated` | `trainer, dest` | Script |
| `CallIfNotDefeated` | `trainer, dest` | Script |
| `ItemVars` | `item, quantity=1` | Script |
| `GoToIfNoItemSpace` | `item, quantity, target` | Script |
| `GoToIfNoItemSpace2` | `item, quantity, target` | Script |
| `GiveItemNoCheck` | `item, quantity` | Script |
| `TakeItemNoCheck` | `item, quantity` | Script |
| `Switch` | `var` | Script |
| `Case` | `value, target` | Script |
| `PhoneCall` | `who, b, c` | Script |
| `SimpleNPCMsg` | `msgid` | Script |
| `ScriptEntry` | `name` | Script |
| `ScriptEntryEnd` | `` | Script |
| `InitScriptEntry_Fixed` | `type, scriptID` | Script |
| `InitScriptEntry_OnFrameTable` | `checksLabel` | Script |
| `InitScriptEntry_OnTransition` | `scriptID` | Script |
| `InitScriptEntry_OnResume` | `scriptID` | Script |
| `InitScriptEntry_OnLoad` | `scriptID` | Script |
| `InitScriptEntryEnd` | `` | Script |
| `InitScriptGoToIfEqual` | `var1, var2, scriptID` | Script |
| `InitScriptFrameTableEnd` | `` | Script |
| `InitScriptEnd` | `` | Script |
| `FaceNorth` | `length=1` | Movement |
| `FaceSouth` | `length=1` | Movement |
| `FaceWest` | `length=1` | Movement |
| `FaceEast` | `length=1` | Movement |
| `WalkSlowerNorth` | `length=1` | Movement |
| `WalkSlowerSouth` | `length=1` | Movement |
| `WalkSlowerWest` | `length=1` | Movement |
| `WalkSlowerEast` | `length=1` | Movement |
| `WalkSlowNorth` | `length=1` | Movement |
| `WalkSlowSouth` | `length=1` | Movement |
| `WalkSlowWest` | `length=1` | Movement |
| `WalkSlowEast` | `length=1` | Movement |
| `WalkNormalNorth` | `length=1` | Movement |
| `WalkNormalSouth` | `length=1` | Movement |
| `WalkNormalWest` | `length=1` | Movement |
| `WalkNormalEast` | `length=1` | Movement |
| `WalkFastNorth` | `length=1` | Movement |
| `WalkFastSouth` | `length=1` | Movement |
| `WalkFastWest` | `length=1` | Movement |
| `WalkFastEast` | `length=1` | Movement |
| `WalkFasterNorth` | `length=1` | Movement |
| `WalkFasterSouth` | `length=1` | Movement |
| `WalkFasterWest` | `length=1` | Movement |
| `WalkFasterEast` | `length=1` | Movement |
| `WalkOnSpotSlowerNorth` | `length=1` | Movement |
| `WalkOnSpotSlowerSouth` | `length=1` | Movement |
| `WalkOnSpotSlowerWest` | `length=1` | Movement |
| `WalkOnSpotSlowerEast` | `length=1` | Movement |
| `WalkOnSpotSlowNorth` | `length=1` | Movement |
| `WalkOnSpotSlowSouth` | `length=1` | Movement |
| `WalkOnSpotSlowWest` | `length=1` | Movement |
| `WalkOnSpotSlowEast` | `length=1` | Movement |
| `WalkOnSpotNormalNorth` | `length=1` | Movement |
| `WalkOnSpotNormalSouth` | `length=1` | Movement |
| `WalkOnSpotNormalWest` | `length=1` | Movement |
| `WalkOnSpotNormalEast` | `length=1` | Movement |
| `WalkOnSpotFastNorth` | `length=1` | Movement |
| `WalkOnSpotFastSouth` | `length=1` | Movement |
| `WalkOnSpotFastWest` | `length=1` | Movement |
| `WalkOnSpotFastEast` | `length=1` | Movement |
| `WalkOnSpotFasterNorth` | `length=1` | Movement |
| `WalkOnSpotFasterSouth` | `length=1` | Movement |
| `WalkOnSpotFasterWest` | `length=1` | Movement |
| `WalkOnSpotFasterEast` | `length=1` | Movement |
| `JumpOnSpotSlowNorth` | `length=1` | Movement |
| `JumpOnSpotSlowSouth` | `length=1` | Movement |
| `JumpOnSpotSlowWest` | `length=1` | Movement |
| `JumpOnSpotSlowEast` | `length=1` | Movement |
| `JumpOnSpotFastNorth` | `length=1` | Movement |
| `JumpOnSpotFastSouth` | `length=1` | Movement |
| `JumpOnSpotFastWest` | `length=1` | Movement |
| `JumpOnSpotFastEast` | `length=1` | Movement |
| `JumpNearFastNorth` | `length=1` | Movement |
| `JumpNearFastSouth` | `length=1` | Movement |
| `JumpNearFastWest` | `length=1` | Movement |
| `JumpNearFastEast` | `length=1` | Movement |
| `JumpFarNorth` | `length=1` | Movement |
| `JumpFarSouth` | `length=1` | Movement |
| `JumpFarWest` | `length=1` | Movement |
| `JumpFarEast` | `length=1` | Movement |
| `Delay1` | `length=1` | Movement |
| `Delay2` | `length=1` | Movement |
| `Delay4` | `length=1` | Movement |
| `Delay8` | `length=1` | Movement |
| `Delay15` | `length=1` | Movement |
| `Delay16` | `length=1` | Movement |
| `Delay32` | `length=1` | Movement |
| `WarpOut` | `length=1` | Movement |
| `WarpIn` | `length=1` | Movement |
| `SetInvisible` | `length=1` | Movement |
| `SetVisible` | `length=1` | Movement |
| `LockDir` | `length=1` | Movement |
| `UnlockDir` | `length=1` | Movement |
| `PauseAnimation` | `length=1` | Movement |
| `ResumeAnimation` | `length=1` | Movement |
| `EmoteExclamationMark` | `length=1` | Movement |
| `WalkSlightlyFastNorth` | `length=1` | Movement |
| `WalkSlightlyFastSouth` | `length=1` | Movement |
| `WalkSlightlyFastWest` | `length=1` | Movement |
| `WalkSlightlyFastEast` | `length=1` | Movement |
| `WalkSlightlyFasterNorth` | `length=1` | Movement |
| `WalkSlightlyFasterSouth` | `length=1` | Movement |
| `WalkSlightlyFasterWest` | `length=1` | Movement |
| `WalkSlightlyFasterEast` | `length=1` | Movement |
| `WalkFastestNorth` | `length=1` | Movement |
| `WalkFastestSouth` | `length=1` | Movement |
| `WalkFastestWest` | `length=1` | Movement |
| `WalkFastestEast` | `length=1` | Movement |
| `RunNorth` | `length=1` | Movement |
| `RunSouth` | `length=1` | Movement |
| `RunWest` | `length=1` | Movement |
| `RunEast` | `length=1` | Movement |
| `JumpNearSlowWest` | `length=1` | Movement |
| `JumpNearSlowEast` | `length=1` | Movement |
| `JumpFartherWest` | `length=1` | Movement |
| `JumpFartherEast` | `length=1` | Movement |
| `WalkEverSoSlightlyFastNorth` | `length=1` | Movement |
| `WalkEverSoSlightlyFastSouth` | `length=1` | Movement |
| `WalkEverSoSlightlyFastWest` | `length=1` | Movement |
| `WalkEverSoSlightlyFastEast` | `length=1` | Movement |
| `NurseJoyBow` | `length=1` | Movement |
| `RevealTrainer` | `length=1` | Movement |
| `PlayerGive` | `length=1` | Movement |
| `EmoteQuestionMark` | `length=1` | Movement |
| `PlayerReceive` | `length=1` | Movement |
| `MoveAction_105` | `length=1` | Movement |
| `MoveAction_106` | `length=1` | Movement |
| `MoveAction_107` | `length=1` | Movement |
| `MoveAction_108` | `length=1` | Movement |
| `MoveAction_109` | `length=1` | Movement |
| `MoveAction_110` | `length=1` | Movement |
| `MoveAction_111` | `length=1` | Movement |
| `MoveAction_112` | `length=1` | Movement |
| `EmoteExclamation2` | `length=1` | Movement |
| `EndMovement` | `` | Movement |
