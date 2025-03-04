declare module 'react-copy-to-clipboard' {
  import { ComponentType, ReactNode } from 'react';

  interface CopyToClipboardProps {
    text: string;
    children: ReactNode;
    onCopy?: (text: string, result: boolean) => void;
  }

  const CopyToClipboard: ComponentType<CopyToClipboardProps>;
  export = CopyToClipboard;
} 