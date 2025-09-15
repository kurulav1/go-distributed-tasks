type Props = { lines: string[] }
export default function JobLog({ lines }: Props) {
  return (
    <div className="grid gap-2 max-w-3xl">
      {lines.slice().reverse().map((l, i) => (
        <pre key={i} className="card overflow-auto whitespace-pre-wrap">{l}</pre>
      ))}
    </div>
  )
}
