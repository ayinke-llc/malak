import { useTheme } from "next-themes";
import BaseSkeleton from "react-loading-skeleton";
import "react-loading-skeleton/dist/skeleton.css";

const Skeleton = ({ count }: { count: number }) => {
  const { theme } = useTheme();

  return (
    <BaseSkeleton
      count={count}
      baseColor={theme === "light" ? "hsl(210 40% 96.1%)" : "hsl(217 32% 17%)"}
      highlightColor={theme === "light" ? "hsl(214.3 31.8% 91.4%)" : "hsl(222 47% 11%)"}
    />
  );
};

export default Skeleton;
