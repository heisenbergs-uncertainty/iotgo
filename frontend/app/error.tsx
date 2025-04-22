'use client';

import Link from 'next/link';
import {useEffect} from 'react';

export default function ErrorPage({
                                    error,
                                  }: {
  error: { message: string; redirect?: string; retryUrl?: string };
}) {
  useEffect(() => {
    document.getElementById('error-timestamp')!.textContent = new Date().toLocaleString();
  }, []);

  return (
      <div className="flex items-center justify-center min-h-[80vh] px-4">
        <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] max-w-lg w-full animate-[var(--animation-fade-in)] border-0 rounded-lg">
          <div className="p-6 text-center">
            <h1 className="text-4xl sm:text-5xl font-bold text-white mb-4">
              Oops! An Error Occurred
            </h1>
            <p className="text-[--color-danger] mb-4">{error.message}</p>
            <p className="text-[--color-text-primary]">
              This error occurred at: <span id="error-timestamp"></span>
            </p>
            <p className="text-[--color-text-primary]">
              If this issue persists, please contact support at{' '}
              <a href="mailto:support@iotgo.com" className="text-[--color-text-accent] hover:underline">
                support@iotgo.com
              </a>
              .
            </p>
            <div className="flex justify-center gap-4 mt-6 flex-wrap">
              {error.redirect && (
                  <Link
                      href={error.redirect}
                      className="bg-[--color-primary] text-white font-semibold py-2 px-4 rounded hover:bg-blue-600 transition-colors"
                  >
                    Go Back
                  </Link>
              )}
              <Link
                  href="/"
                  className="bg-[--color-secondary] text-white font-semibold py-2 px-4 rounded hover:bg-gray-600 transition-colors"
              >
                Return to Homepage
              </Link>
              {error.retryUrl && (
                  <Link
                      href={error.retryUrl}
                      className="border border-[--color-primary] text-[--color-primary] font-semibold py-2 px-4 rounded hover:bg-[--color-primary] hover:text-white transition-colors"
                  >
                    Retry Action
                  </Link>
              )}
            </div>
          </div>
        </div>
      </div>
  );
}