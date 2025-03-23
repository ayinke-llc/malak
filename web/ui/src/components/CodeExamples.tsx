import { useEffect, useState } from 'react';
import * as shiki from 'shiki';
import { Card } from './ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Button } from './ui/button';
import { RiCheckLine, RiFileCopyLine } from '@remixicon/react';
import { toast } from 'sonner';

interface CodeExample {
  language: string;
  code: string;
  label: string;
}

interface CodeExamplesProps {
  examples: CodeExample[];
}

export function CodeExamples({ examples }: CodeExamplesProps) {
  const [highlightedCode, setHighlightedCode] = useState<Record<string, string>>({});
  const [copiedLanguage, setCopiedLanguage] = useState<string | null>(null);

  useEffect(() => {
    const highlightCode = async () => {
      const highlighter = await shiki.createHighlighter({
        themes: ['github-light'],
        langs: ['typescript', 'python', 'go', 'shell'],
      });

      const highlighted: Record<string, string> = {};

      for (const example of examples) {
        highlighted[example.language] = highlighter.codeToHtml(example.code, {
          lang: example.language,
          theme: 'github-light',
        });
      }

      setHighlightedCode(highlighted);
    };

    highlightCode();
  }, [examples]);

  const copyCode = async (language: string, code: string) => {
    try {
      await navigator.clipboard.writeText(code);
      setCopiedLanguage(language);
      toast.success('Code copied to clipboard');
      setTimeout(() => setCopiedLanguage(null), 2000);
    } catch (err) {
      toast.error('Failed to copy code');
    }
  };

  return (
    <Card className="w-full">
      <Tabs defaultValue={examples[0]?.language} className="w-full">
        <TabsList className="w-full">
          {examples.map((example) => (
            <TabsTrigger
              key={example.language}
              value={example.language}
              className="flex-1"
            >
              {example.label}
            </TabsTrigger>
          ))}
        </TabsList>
        {examples.map((example) => (
          <TabsContent
            key={example.language}
            value={example.language}
            className="mt-0 relative group"
          >
            <div className="absolute right-2 top-2 z-10">
              <Button
                variant="ghost"
                size="sm"
                className="gap-2"
                onClick={() => copyCode(example.language, example.code)}
              >
                {copiedLanguage === example.language ? (
                  <>
                    <RiCheckLine className="h-4 w-4" />
                    Copied!
                  </>
                ) : (
                  <>
                    <RiFileCopyLine className="h-4 w-4" />
                    Copy
                  </>
                )}
              </Button>
            </div>
            <div
              className="rounded-md overflow-hidden bg-[#ffffff] border"
              dangerouslySetInnerHTML={{ __html: highlightedCode[example.language] || '' }}
            />
          </TabsContent>
        ))}
      </Tabs>
    </Card>
  );
} 
