import styles from './Page.module.css'

interface ICellProps {
  children: React.ReactNode
}

const Cell = ({ children }: ICellProps) => {
  return (
    <td>
      <div className={styles.cellScroll}>
        {children}
      </div>
    </td>
  )
}

const Page = () => {
  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>Timestamp</th>
            <th>ID</th>
            <th>Module</th>
            <th>Action</th>
            <th>Target</th>
            <th>Status</th>
            <th>Duration</th>
            <th>Confidence</th>
            <th>Summary</th>
            <th>Parent</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <Cell>2025-10-29 15:42:10</Cell>
            <Cell>run-42 / step-007</Cell>
            <Cell>LLM</Cell>
            <Cell>Cluster anomalies</Cell>
            <Cell>api.acme.shop /v1/orders</Cell>
            <Cell>OK</Cell>
            <Cell>412ms</Cell>
            <Cell>0.71</Cell>
            <Cell>Possible IDOR on order detail.</Cell>
            <Cell>-1</Cell>
          </tr>
        </tbody>
      </table>
    </div>
  )
}

export default Page