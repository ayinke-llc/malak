import { usePathname } from "next/navigation";

export default function MobileSidebar() {
  const pathname = usePathname();

  const isActive = (itemHref: string) => {
    return pathname === itemHref || pathname.startsWith(itemHref);
  };

  return null
}

