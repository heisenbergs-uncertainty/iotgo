'use client';

import Link from 'next/link';
import {usePathname} from 'next/navigation';
import {useCallback, useEffect, useRef, useState} from 'react';
import {useTheme} from '@/contexts/ThemeContext';
import debounce from 'lodash.debounce';
import {IoTGoLogo} from "@/components/visual/IoTGoLogo";
import {User} from "@/types/api";

const Navbar = () => {
  const user: User = {Id: 1, Username: 'admin', Role: 'admin'}; // Replace with actual user data from context
  const pathname = usePathname();
  const [isOpen, setIsOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const { theme, toggleTheme } = useTheme();
  const menuRef = useRef<HTMLDivElement>(null);

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setUserMenuOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Close mobile menu when route changes
  useEffect(() => {
    setIsOpen(false);
  }, [pathname]);

  // Debounced scroll handler
  const handleScroll = useCallback(
      debounce(() => {
        setScrolled(window.scrollY > 10);
      }, 100),
      []
  );

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [handleScroll]);



  const handleLogout = async () => {}
  //if (isLoading) {
   // return <nav className="fixed top-0 w-full bg-[--color-background-secondary] h-16"></nav>;
  //}

  return (
      <nav
          className={`fixed top-0 w-full bg-[--color-background-secondary] z-50 transition-all duration-300 ${
              scrolled ? 'shadow-[0_4px_12px_var(--color-card-shadow)]' : ''
          }`}
      >
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link
                href="/"
                className="flex items-center hover:opacity-90 transition-all"
                aria-label="IoTGo Home"
            >
              <IoTGoLogo className="h-9" />
              <span className="sr-only">IoTGo</span>
            </Link>
            {/* Mobile menu button */}
            <div className="lg:hidden flex items-center">
              {/* Theme Toggle Button (Mobile) */}
              <button
                  onClick={toggleTheme}
                  aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
                  className="p-2 mr-2 rounded-full hover:bg-gray-700/30 text-[--color-text-primary] hover:text-[--color-text-accent] transition-colors flex items-center"
              >
                {theme === 'dark' ? (
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <circle cx="12" cy="12" r="5"/>
                      <line x1="12" y1="1" x2="12" y2="3"/>
                      <line x1="12" y1="21" x2="12" y2="23"/>
                      <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
                      <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
                      <line x1="1" y1="12" x2="3" y2="12"/>
                      <line x1="21" y1="12" x2="23" y2="12"/>
                      <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
                      <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
                    </svg>
                ) : (
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
                    </svg>
                )}
              </button>

              <button
                  className="text-[--color-text-primary] p-2 rounded-md hover:bg-gray-700/30 transition-colors"
                  onClick={() => setIsOpen((prev) => !prev)}
                  aria-label={isOpen ? 'Close menu' : 'Open menu'}
                  aria-expanded={isOpen}
              >
                <svg
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                    className="transition-transform duration-300"
                    style={{ transform: isOpen ? 'rotate(90deg)' : 'rotate(0)' }}
                >
                  {isOpen ? (
                      <path
                          d="M18 6L6 18M6 6L18 18"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                      />
                  ) : (
                      <path
                          d="M4 6H20M4 12H20M4 18H20"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                      />
                  )}
                </svg>
              </button>
            </div>

            {/* Navigation Links and Actions (Desktop) */}
            <div className="hidden lg:flex items-center flex-1 justify-between ml-8">
              {/* Nav Links */}
              <ul className="flex items-center space-x-6">
                <li>
                  <Link
                      href="/devices"
                      className={`flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-1 relative group ${
                          pathname.startsWith('/devices') ? 'text-[--color-text-accent]' : ''
                      }`}
                  >
                    Devices
                    <span className={`absolute -bottom-1 left-0 w-0 h-0.5 bg-[--color-text-accent] transition-all duration-300 group-hover:w-full ${
                        pathname.startsWith('/devices') ? 'w-full' : ''
                    }`}></span>
                  </Link>
                </li>

                {user?.Role === 'admin' && (
                    <li>
                      <Link
                          href="/admin"
                          className={`flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-1 relative group ${
                              pathname === '/admin' ? 'text-[--color-text-accent]' : ''
                          }`}
                      >
                        Admin Dashboard
                        <span className={`absolute -bottom-1 left-0 w-0 h-0.5 bg-[--color-text-accent] transition-all duration-300 group-hover:w-full ${
                            pathname === '/admin' ? 'w-full' : ''
                        }`}></span>
                      </Link>
                    </li>
                )}
              </ul>

              {/* Right side items container */}
              <div className="flex items-center space-x-4">
                {/* User Menu */}
                {user ? (
                    <div className="relative group" ref={menuRef}>
                      <button
                          onClick={() => setUserMenuOpen(!userMenuOpen)}
                          className="flex items-center space-x-2 text-[--color-text-primary] font-medium py-2 px-3 rounded-md hover:bg-gray-700/30 transition-colors"
                          aria-expanded={userMenuOpen}
                          aria-haspopup="true"
                      >
                        <span>{user.Username}</span>
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="2"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            className={`transition-transform duration-200 ${userMenuOpen ? 'rotate-180' : ''}`}
                        >
                          <polyline points="6 9 12 15 18 9"></polyline>
                        </svg>
                      </button>

                      {/* User Dropdown Menu */}
                      {userMenuOpen && (
                          <div className="absolute right-0 mt-2 w-48 bg-[--color-background-secondary] rounded-lg shadow-lg border border-gray-700/30 py-1 z-50">
                            <Link
                                href="/profile"
                                className="block px-4 py-2 text-[--color-text-primary] hover:bg-gray-700/30 hover:text-[--color-text-accent] transition-colors w-full text-left"
                                onClick={() => setUserMenuOpen(false)}
                            >
                              Profile
                            </Link>
                            <button
                                onClick={handleLogout}
                                className="block px-4 py-2 text-[--color-text-primary] hover:bg-gray-700/30 hover:text-[--color-text-accent] transition-colors w-full text-left"
                            >
                              Logout
                            </button>
                          </div>
                      )}
                    </div>
                ) : (
                    <Link
                        href="/login"
                        className="flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-3 rounded-md hover:bg-gray-700/30"
                    >
                      <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="18"
                          height="18"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          className="mr-2"
                      >
                        <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path>
                        <polyline points="10 17 15 12 10 7"></polyline>
                        <line x1="15" y1="12" x2="3" y2="12"></line>
                      </svg>
                      Login
                    </Link>
                )}

                {/* Theme Toggle Button (Desktop) */}
                <button
                    onClick={toggleTheme}
                    aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
                    className="p-2 rounded-full hover:bg-gray-700/30 text-[--color-text-primary] hover:text-[--color-text-accent] transition-colors"
                >
                  {theme === 'dark' ? (
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <circle cx="12" cy="12" r="5"/>
                        <line x1="12" y1="1" x2="12" y2="3"/>
                        <line x1="12" y1="21" x2="12" y2="23"/>
                        <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
                        <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
                        <line x1="1" y1="12" x2="3" y2="12"/>
                        <line x1="21" y1="12" x2="23" y2="12"/>
                        <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
                        <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
                      </svg>
                  ) : (
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
                      </svg>
                  )}
                </button>
              </div>
            </div>

            {/* Mobile Navigation Menu */}
            {isOpen && (
                <div className="absolute top-16 left-0 w-full bg-[--color-background-secondary] shadow-[0_4px_12px_var(--color-card-shadow)] border-t border-gray-700/50 lg:hidden">
                  <ul className="flex flex-col space-y-2 p-4">
                    <li>
                      <Link
                          href="/devices"
                          className={`flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-1 relative group ${
                              pathname.startsWith('/devices') ? 'text-[--color-text-accent]' : ''
                          }`}
                      >
                        Devices
                        <span className={`absolute -bottom-1 left-0 w-0 h-0.5 bg-[--color-text-accent] transition-all duration-300 group-hover:w-full ${
                            pathname.startsWith('/devices') ? 'w-full' : ''
                        }`}></span>
                      </Link>
                    </li>

                    {user?.Role === 'admin' && (
                        <li>
                          <Link
                              href="/admin"
                              className={`flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-1 relative group ${
                                  pathname === '/admin' ? 'text-[--color-text-accent]' : ''
                              }`}
                          >
                            Admin Dashboard
                            <span className={`absolute -bottom-1 left-0 w-0 h-0.5 bg-[--color-text-accent] transition-all duration-300 group-hover:w-full ${
                                pathname === '/admin' ? 'w-full' : ''
                            }`}></span>
                          </Link>
                        </li>
                    )}
                  </ul>

                  {/* Mobile User Menu */}
                  <div className="flex flex-col p-4 border-t border-gray-700/30">
                    {user ? (
                        <>
                          <Link
                              href="/profile"
                              className="flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-1"
                          >
                            Profile
                          </Link>
                          <button
                              onClick={handleLogout}
                              className="flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-1 text-left"
                          >
                            Logout
                          </button>
                        </>
                    ) : (
                        <Link
                            href="/login"
                            className="flex items-center text-[--color-text-primary] hover:text-[--color-text-accent] font-medium transition-colors py-2 px-3 rounded-md hover:bg-gray-700/30"
                        >
                          <svg
                              xmlns="http://www.w3.org/2000/svg"
                              width="18"
                              height="18"
                              viewBox="0 0 24 24"
                              fill="none"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              className="mr-2"
                          >
                            <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path>
                            <polyline points="10 17 15 12 10 7"></polyline>
                            <line x1="15" y1="12" x2="3" y2="12"></line>
                          </svg>
                          Login
                        </Link>
                    )}
                  </div>
                </div>
            )}
          </div>
        </div>
      </nav>
  );
};

export default Navbar;