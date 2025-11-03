export type AbilityFlag = number;

export const MULTISHOT_ABILITY_FLAG: AbilityFlag = 1 << 1;
export const WIDE_BEAM_ABILITY_FLAG: AbilityFlag = 1 << 2;
export const SHIELD_ABILITY_FLAG: AbilityFlag = 1 << 3;

/**
 * Reports whether an ability is active.
 * @param flags current active flags
 * @param abilityFlag the target ability
 */
export function isAbilityActive(
  flags: AbilityFlag,
  abilityFlag: AbilityFlag,
): boolean {
  return (flags & abilityFlag) != 0;
}

/**
 * Converts an ability to its name (to index spritesheet).
 * Expects exactly one active flag. Returns null otherwise.
 * @param flag current ability flag
 */
export function toAbilityName(flag: AbilityFlag): string | null {
  switch (flag) {
  case MULTISHOT_ABILITY_FLAG:
    return "multishot";
  case WIDE_BEAM_ABILITY_FLAG:
    return "wide";
  case SHIELD_ABILITY_FLAG:
    return "shield";
  default:
    return null;
  }
}
