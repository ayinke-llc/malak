import { useTheme } from "next-themes";
import BaseSkeleton from "react-loading-skeleton";
import "react-loading-skeleton/dist/skeleton.css";

const Skeleton = ({ count }: { count: number }) => {

  const { theme } = useTheme()

  return <BaseSkeleton
    count={count}
    baseColor={theme === "light" ? "black" : "white"}
    highlightColor="#595352" />;
};

export default Skeleton;
