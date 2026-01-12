import './HealthStatus.css'

function HealthStatus({ health }) {
  if (!health) {
    return (
      <div className="health-status unknown">
        <span className="status-indicator">⏳</span>
        <span>Checking backend status...</span>
      </div>
    )
  }

  const isNodeHealthy = health.status === 'ok'
  const goBackend = health.goBackend
  const isGoHealthy = goBackend?.status === 'ok'
  const allHealthy = isNodeHealthy && isGoHealthy

  return (
    <div className={`health-status ${allHealthy ? 'healthy' : 'unhealthy'}`}>
      <div className="health-row">
        <span className="status-indicator">
          {isNodeHealthy ? '✅' : '❌'}
        </span>
        <span>Node.js API Gateway: {health.message || 'running'}</span>
      </div>
      
      {goBackend && (
        <div className="health-row">
          <span className="status-indicator">
            {isGoHealthy ? '✅' : '❌'}
          </span>
          <span>
            Go Backend: {goBackend.message || 'running'}
            {goBackend.version && ` (v${goBackend.version})`}
            {goBackend.uptime && ` - uptime: ${goBackend.uptime}`}
          </span>
        </div>
      )}

      {goBackend?.checks && (
        <div className="health-checks">
          {Object.entries(goBackend.checks).map(([key, value]) => (
            <span key={key} className="check-item">
              {value === 'ok' ? '✓' : '✗'} {key}
            </span>
          ))}
        </div>
      )}
    </div>
  )
}

export default HealthStatus
