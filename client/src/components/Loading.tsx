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
        className={`${sizeClasses[size]} animate-spin rounded-full border-4 border-gray-200 border-t-blue-600 dark:border-gray-700 dark:border-t-blue-400`}
      />
    </div>
  );
}

// Loading skeleton for product cards
export function ProductCardSkeleton() {
  return (
    <div className="flex flex-col overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm dark:border-gray-800 dark:bg-gray-900">
      <div className="aspect-square w-full animate-pulse bg-gray-200 dark:bg-gray-800" />
      <div className="flex flex-1 flex-col p-4">
        <div className="mb-2 h-5 w-3/4 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
        <div className="mb-3 h-4 w-full animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
        <div className="mt-auto flex items-center justify-between">
          <div className="h-6 w-24 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
          <div className="h-4 w-16 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
        </div>
      </div>
    </div>
  );
}

