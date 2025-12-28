export default function Footer() {
  return (
    <footer className="border-t border-neutral-200 bg-white">
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="text-center">
          <p className="text-sm text-neutral-500">
            &copy; {new Date().getFullYear()} Ecommerce Store. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}

