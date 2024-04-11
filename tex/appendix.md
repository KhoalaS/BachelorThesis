\clearpage
\appendix

# Appendix

$$
\begin{table}[h]
\centering
\begin{subtable}[b]{\textwidth}
\makebox[\textwidth][c]{
\begin{tabular}{lrrrrrrrrrrrrrr}
\toprule
 & ratio & Tiny & VD & ED & Small & Tri & ETri & AVD & ADVD & SED2 & SED2* & F3 & $|C|$ & est. opt \\
\midrule
mean & 1.9196 & 59.99 & 345.23 & 2.40 & 0.46 & 0.04 & 47.26 & 0.03 & 2.02 & 44.47 & 9.16 & 44.33 & 592.46 & 308.68 \\
std & 0.0273 & 5.86 & 9.71 & 1.56 & 0.63 & 0.19 & 5.27 & 0.18 & 1.35 & 3.93 & 2.67 & 6.54 & 8.59 & 3.51 \\
min & 1.8344 & 43 & 316 & 0 & 0 & 0 & 29 & 0 & 0 & 31 & 2 & 26 & 564 & 298 \\
median & 1.9199 & 60 & 345 & 2 & 0 & 0 & 47 & 0 & 2 & 44 & 9 & 45 & 592 & 309 \\
max & 2.0364 & 79 & 375 & 8 & 3 & 1 & 61 & 2 & 7 & 57 & 19 & 70 & 622 & 320 \\
\bottomrule
\end{tabular}
}
\caption{random edge in F3 rule}
\end{subtable}
\newline
\vspace{4mm}
\newline
\begin{subtable}[b]{\textwidth}
\makebox[\textwidth][c]{
\begin{tabular}{lrrrrrrrrrrrrrr}
\toprule
 & ratio & Tiny & VD & ED & Small & Tri & ETri & AVD & ADVD & SED2 & SED2* & F3 & $|C|$ & est. opt \\
\midrule
mean & 1.8690 & 58.28 & 349.60 & 2.33 & 0.47 & 0.05 & 63.68 & 0.03 & 2.01 & 42.57 & 8.59 & 25.51 & 590.66 & 316.06 \\
std & 0.0225 & 5.93 & 9.64 & 1.57 & 0.62 & 0.21 & 4.10 & 0.16 & 1.39 & 4.06 & 2.70 & 3.73 & 8.98 & 3.3400 \\
min & 1.7987 & 41 & 321 & 0 & 0 & 0 & 50 & 0 & 0 & 29 & 1 & 16 & 562 & 305 \\
median & 1.8675 & 58 & 350 & 2 & 0 & 0 & 64 & 0 & 2 & 43 & 9 & 25 & 590 & 316 \\
max & 1.9525 & 75 & 382 & 10 & 3 & 1 & 77 & 1 & 9 & 56 & 17 & 40 & 617 & 328 \\
\bottomrule
\end{tabular}
}
\caption{F3 low degree rule}
\end{subtable}
\caption{Results for 100 random 3-uniform ER hypergraphs with an edge to vertex ratio of 3. Data was collected 10 times per graph to account for run-to-run variance.}
\label{er3_evr3_stats}
\end{table}
$$

$$
\begin{figure}[h]
    \centering
    \includegraphics[width=\textwidth]{./img/er3_evr_f3_rules}
    \caption{Mean number of rule executions for 100 random 3-uniform ER hypergraphs with an edge to vertex ratio of 3. Data was collected 10 times per graph to account for run to run variance.}
    \label{fig:er3_evr3_rules}
\end{figure}
$$

$$
\begin{table}[h]
\centering
\begin{tabular}{lrrrr}
\toprule
 & ratio & $|C|$ & est. opt & time \\
\midrule
mean & 1.3142 & 80927.77 & 61579.82 & 168 sec\\
std & 0.0006 & 55.84 & 32.98 & 3 sec\\
min & 1.3128 & 80786 & 61511 & 161 sec\\
median & 1.3142 & 80926.50 & 61579.50 & 168 sec\\
max & 1.3157 & 81071 & 61654 & 176 sec\\
\bottomrule
\end{tabular}
\caption{Results for \textsc{Triangle Vertex Deletion} on Amazon product co-purchasing graph using rule strategy 3; F3 rule selects a random size three edge;  $n=100$}
\label{amzn_str_3_f3rand}
\end{table}
$$

$$
\begin{figure}[h]
    \centering
    \includegraphics[width=\textwidth]{./img/flamegraph.png}
    \caption{Flamegraph of pprof CPU performance profile. Graph instance was a dense 3-uniform ER hypergraph with 1000 vertices. }
    \label{flamegraph}
\end{figure}
$$

$$
\begin{table}[h]
    \begin{subtable}[b]{0.45\textwidth}
        \input{"out/rome_cvd_lphs_glpk"}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \input{"out/rome_cvd_lpsc_glpk"}
    \end{subtable}
    \newline
    \vspace{4mm}
    \newline
    \begin{subtable}[b]{0.45\textwidth}
        \input{"out/rome_cvd_lphs_clp"}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \input{"out/rome_cvd_lpsc_clp"}
    \end{subtable}
    \newline
    \vspace{4mm}
    \newline
    \begin{subtable}[b]{0.45\textwidth}
        \input{"out/rome_cvd_lphs_highs"}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \input{"out/rome_cvd_lpsc_highs"}
    \end{subtable}
    \newline
    \vspace{4mm}
    \newline
    \begin{subtable}[t]{0.45\textwidth}
        \input{"out/rome_cvd_lphs_highs_ipm"}
    \end{subtable}
    \hfill
    \begin{subtable}[t]{0.45\textwidth}
        \input{"out/rome_cvd_lpsc_highs_ipm"}
    \end{subtable}
    \caption{Results for \textsc{Cluster Vertex Deletion} with LP rounding based algorithms on graphs from the Rome Graphs collection. Same method as previous CVD benchmark. Assume simplex method if not specifed otherwise.}
    \label{A:cvd_rome_all}
\end{table}
$$