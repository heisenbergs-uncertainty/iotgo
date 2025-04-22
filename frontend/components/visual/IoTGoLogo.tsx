// Logo Component for the Navbar
import {useTheme} from "@/contexts/ThemeContext";

export const IoTGoLogo = ({ className = "" }: { className?: string }) => {
    const { theme } = useTheme();

    return (
        <svg
            width="120"
            height="36"
            viewBox="0 0 120 36"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className={className}
        >
            {/* Circuit/Connection Elements */}
            <path
                d="M12 10C16.4183 10 20 13.5817 20 18C20 22.4183 16.4183 26 12 26"
                stroke={theme === 'dark' ? '#00C4B4' : '#007BFF'}
                strokeWidth="2.5"
                strokeLinecap="round"
            />
            <circle
                cx="12"
                cy="18"
                r="3"
                fill={theme === 'dark' ? '#00C4B4' : '#007BFF'}
            />
            <path
                d="M24 10V26"
                stroke={theme === 'dark' ? '#00C4B4' : '#007BFF'}
                strokeWidth="2.5"
                strokeLinecap="round"
                strokeDasharray="1 3"
            />

            {/* Text Elements */}
            <text
                x="32"
                y="23"
                fontFamily="Inter, sans-serif"
                fontSize="20"
                fontWeight="700"
                fill={theme === 'dark' ? '#00C4B4' : '#007BFF'}
            >
                IoTGo
            </text>

            {/* Dots representing IoT network */}
            <circle cx="84" cy="12" r="2" fill={theme === 'dark' ? '#E0E0E0' : '#333333'} />
            <circle cx="92" cy="16" r="1.5" fill={theme === 'dark' ? '#E0E0E0' : '#333333'} />
            <circle cx="100" cy="14" r="1" fill={theme === 'dark' ? '#E0E0E0' : '#333333'} />
        </svg>
    );
};
