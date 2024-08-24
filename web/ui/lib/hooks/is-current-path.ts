export function useIsCurrentPath(path: string): boolean {
  if (typeof window !== 'undefined') {
    return window.location.pathname === path;
  }
  return false;
}
