export type AbilityFlag = number;

export const MULTISHOT_ABILITY_FLAG: AbilityFlag = 1 << 1;
export const WIDE_BEAM_ABILITY_FLAG: AbilityFlag = 1 << 2;

export function isAbilityActive(flags: AbilityFlag, abilityFlag: AbilityFlag): boolean {
  return (flags & abilityFlag) != 0;
}
