function detectDesktopOS(): "macos" | "windows" | "linux" | "unknown" {
  if (typeof navigator === "undefined") return "unknown";
  const ua = navigator.userAgent.toLowerCase();
  const platform = (navigator.platform ?? "").toLowerCase();
  if (platform.includes("mac") || ua.includes("mac os")) return "macos";
  if (platform.includes("win") || ua.includes("windows")) return "windows";
  if (platform.includes("linux") || ua.includes("linux")) return "linux";
  return "unknown";
}

function detectNativeMobile(): boolean {
  if (typeof navigator === "undefined") return false;
  const ua = navigator.userAgent.toLowerCase();
  return /android|iphone|ipad|ipod/.test(ua);
}

class PlatformStore {
  viewportWidth = $state(typeof window !== "undefined" ? window.innerWidth : 1024);
  isTouch = $state(false);
  isNativeMobile = $state(detectNativeMobile());
  desktopOS = $state<"macos" | "windows" | "linux" | "unknown">(detectDesktopOS());
  nativeEffect = $state<"None">("None");

  isMobile = $derived(this.viewportWidth < 640);
  isTablet = $derived(this.viewportWidth >= 640 && this.viewportWidth < 1024);
  isDesktop = $derived(this.viewportWidth >= 1024);
  hasWindowChrome = $derived(!this.isNativeMobile);
  hasNativeBlur = $derived(false);

  constructor() {
    if (typeof window === "undefined") return;

    const onResize = () => {
      this.viewportWidth = window.innerWidth;
    };
    window.addEventListener("resize", onResize);

    const mq = window.matchMedia("(pointer: coarse)");
    this.isTouch = mq.matches;
    mq.addEventListener("change", (e) => {
      this.isTouch = e.matches;
    });
  }
}

export const platformStore = new PlatformStore();