import "./globals.css";
import ClientLayout from "@/components/ClientLayout";

export const metadata = {
  title: "IoTGo Platform",
  description:
    "Comprehensive solution for IoT device management and monitoring",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {

  return (
    <html lang="en" className="dark">
      <body>
        <ClientLayout>{children}</ClientLayout>
      </body>
    </html>
  );
}
