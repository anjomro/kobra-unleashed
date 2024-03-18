import { DateTime, DateTimeFormatOptions } from 'luxon';

export const convertSize = (size: number) => {
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];

  let i = 0;
  while (size >= 1024) {
    size /= 1024;
    i++;
  }

  return `${size.toFixed()} ${units[i]}`;
};

export const convertTimestamp = (timestamp: string) => {
  // try {
  //   // Get browser locale
  //   const locale = navigator.language as string;

  //   // Convert timestamp to date object
  //   const dateTime = DateTime.fromMillis(timestamp);

  //   // Extract time using browser locale
  //   const time = dateTime.toLocaleString(
  //     {
  //       hour: 'numeric',
  //       minute: 'numeric',
  //       second: 'numeric',
  //     },
  //     {
  //       locale,
  //     }
  //   );

  //   // Format date example: "April 1, 2022"
  //   const format: DateTimeFormatOptions = {
  //     month: 'long',
  //     day: 'numeric',
  //     year: 'numeric',
  //   };
  //   // Combine formatted date and time with a space
  //   return `${dateTime.toLocaleString(format)} ${time}`;
  // } catch (error) {
  //   // Handle errors
  //   console.error('Error formatting timestamp:', error);
  //   return 'Invalid timestamp';
  // }

  // Do above but: 2024-03-11T03:34:40+08:00
  const dateTime = DateTime.fromISO(timestamp);

  // Convert to desired format
  // return `${dateTime.toLocaleString(format)} ${time}`;

  // Get browser locale
  const locale = navigator.language as string;

  // Extract time using browser locale
  const time = dateTime.toLocaleString(
    {
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
    },
    {
      locale,
    }
  );

  // Format date example: "April 1, 2022"
  const format: DateTimeFormatOptions = {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  };
  // Combine formatted date and time with a space
  return `${dateTime.toLocaleString(format)} ${time}`;
};
