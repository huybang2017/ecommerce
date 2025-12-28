interface ErrorProps {
  message: string;
  onRetry?: () => void;
  className?: string;
}

export default function Error({ message, onRetry, className = '' }: ErrorProps) {
  return (
    <div className={`flex flex-col items-center justify-center p-12 ${className}`}>
      <div className="mb-6 text-red-500">
        <svg
          className="h-16 w-16"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
      </div>
      <p className="mb-6 text-center text-lg font-medium text-neutral-600">
        {message}
      </p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="rounded-lg bg-neutral-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-neutral-800"
        >
          Retry
        </button>
      )}
    </div>
  );
}

