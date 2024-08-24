"use client"

const ComingSoon = () => {
  return (
    <div style={{
      height: '100vh', // Changed from minHeight to height
      width: '100%', // Added to ensure full 
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
      color: 'white',
      fontFamily: 'Arial, sans-serif'
    }}>
      <main style={{ textAlign: 'center', padding: '20px' }}>
        <h1 style={{
          fontSize: '3rem',
          marginBottom: '20px',
          position: 'relative',
          display: 'inline-block'
        }}>
          New Feature Coming Soon
          <span style={{
            position: 'absolute',
            top: '0',
            left: '0',
            right: '0',
            bottom: '0',
            background: 'rgba(255,255,255,0.1)',
            transform: 'skew(5deg)',
            zIndex: -1
          }}></span>
        </h1>

        <p style={{ fontSize: '1.2rem', marginBottom: '40px', maxWidth: '600px' }}>
          We're working on an exciting new feature to enhance your dashboard experience.
          Stay tuned for the big reveal!
        </p>

        <div style={{
          position: 'relative',
          width: '200px',
          height: '200px',
          margin: '0 auto'
        }}>
          <div style={{
            position: 'absolute',
            width: '100%',
            height: '100%',
            border: '4px solid rgba(255,255,255,0.7)',
            borderRadius: '50%',
            borderLeftColor: 'transparent',
            borderBottomColor: 'transparent',
            animation: 'spin 3s linear infinite'
          }}></div>
          <div style={{
            position: 'absolute',
            width: '70%',
            height: '70%',
            top: '15%',
            left: '15%',
            border: '4px solid rgba(255,255,255,0.5)',
            borderRadius: '50%',
            borderTopColor: 'transparent',
            borderRightColor: 'transparent',
            animation: 'spin 2s linear infinite reverse'
          }}></div>
          <div style={{
            position: 'absolute',
            width: '40%',
            height: '40%',
            top: '30%',
            left: '30%',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            fontSize: '2rem'
          }}>
            🚀
          </div>
        </div>

        <p style={{ fontSize: '1rem', marginTop: '40px', opacity: 0.8 }}>
          We appreciate your patience as we work to improve your experience.
        </p>
      </main>

      <style jsx>{`
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
};

export default ComingSoon;
