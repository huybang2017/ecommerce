interface LoadingProps {
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export default function Loading({ size = 'md', className = '' }: LoadingProps) {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  };

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <div
        className={`${sizeClasses[size]} animate-spin rounded-full border-4 border-neutral-200 border-t-neutral-900`}
      />
    </div>
  );
}

// Loading skeleton for product cards
export function ProductCardSkeleton() {
  return (
    <div className="flex flex-col overflow-hidden rounded-xl border border-neutral-200 bg-white">
      <div className="aspect-square w-full animate-pulse bg-neutral-100" />
      <div className="flex flex-1 flex-col p-5">
        <div className="mb-2 h-5 w-3/4 animate-pulse rounded bg-neutral-100" />
        <div className="mb-4 h-4 w-full animate-pulse rounded bg-neutral-100" />
        <div className="mt-auto flex items-center justify-between">
          <div className="h-6 w-24 animate-pulse rounded bg-neutral-100" />
          <div className="h-4 w-16 animate-pulse rounded bg-neutral-100" />
        </div>
      </div>
    </div>
  );
}

