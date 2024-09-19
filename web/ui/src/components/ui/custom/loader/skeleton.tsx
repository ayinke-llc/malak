import BaseSkeleton from 'react-loading-skeleton';
import 'react-loading-skeleton/dist/skeleton.css'

const Skeleton = ({ count }: { count: number }) => {
  return (
    <BaseSkeleton count={count} />
  )
}

export default Skeleton;
