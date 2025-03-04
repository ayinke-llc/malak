'use client';

import { Component } from 'react';
import * as Sentry from '@sentry/nextjs';

interface ErrorBoundaryProps {
  children: React.ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  err: Error
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state = { hasError: false, err: new Error() };

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, err: error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    Sentry.captureException(error);
  }

  render() {
    if (this.state.hasError) {
      Sentry.captureException(this.state.err);
      return <h1>Something went wrong.</h1>;
    }

    return this.props.children;
  }
} 
