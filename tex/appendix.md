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
 & est. ratio & Tiny & VD & ED & Small & Tri & ETri & AVD & ADVD & SED2 & SED2* & F3 & $|C|$ & est. opt \\
\midrule
mean & 1.9196 & 59.99 & 345.23 & 2.40 & 0.46 & 0.04 & 47.26 & 0.03 & 2.02 & 44.47 & 9.16 & 44.33 & 592.46 & 308.68 \\
std & 0.0273 & 5.86 & 9.71 & 1.56 & 0.63 & 0.19 & 5.27 & 0.18 & 1.35 & 3.93 & 2.67 & 6.54 & 8.59 & 3.51 \\
min & 1.8344 & 43 & 316 & 0 & 0 & 0 & 29 & 0 & 0 & 31 & 2 & 26 & 564 & 298 \\
median & 1.9199 & 60 & 345 & 2 & 0 & 0 & 47 & 0 & 2 & 44 & 9 & 45 & 592 & 309 \\
max & 2.0364 & 79 & 375 & 8 & 3 & 1 & 61 & 2 & 7 & 57 & 19 & 70 & 622 & 320 \\
\bottomrule
\end{tabular}
}
\caption{random edge in fallback rule}
\end{subtable}
\newline
\vspace{4mm}
\newline
\begin{subtable}[b]{\textwidth}
\makebox[\textwidth][c]{
\begin{tabular}{lrrrrrrrrrrrrrr}
\toprule
 & est. ratio & Tiny & VD & ED & Small & Tri & ETri & AVD & ADVD & SED2 & SED2* & F3 & $|C|$ & est. opt \\
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
\caption{Results for 100 random 3-uniform ER hypergraphs with 1000 vertices and about 3000 edges. Data was collected 10 times per graph to account for run-to-run variance.}
\label{er3_evr3_stats}
\end{table}
$$

\input{out/f3_logs}

\clearpage 

$$
\begin{figure}[h]
    \centering
    \includegraphics[width=\textwidth]{./img/er3_evr_f3_rules}
    \caption{Mean number of rule executions for 100 random 3-uniform ER hypergraphs with an edge to vertex ratio of 3. Data was collected 10 times per graph to account for run to run variance.}
    \label{fig:er3_evr3_rules}
\end{figure}
$$


$$
\begin{figure}[h]
    \centering
    \begin{subfigure}[b]{\textwidth}
        \includegraphics[width=\textwidth]{img/flame_fr_base.png}
        \caption{3-unifrom ER graph with 1000 vertices and 20000 edges, using rule strategy 2}
        \label{A:flame_base}
   \end{subfigure}
    \newline
    \vspace{4mm}
    \newline
    \begin{subfigure}[b]{\textwidth}
        \includegraphics[width=\textwidth]{img/flame_fr_dblp.png}
        \caption{\textsc{Triangle Vertex Deletion} instance from DBLP coauthor graph, using rule strategy 3}
        \label{A:flame_dblp}
    \end{subfigure}
    \caption{Flamegraph of pprof CPU performance profile. Sections marked with a red box are related to frontier expansion.}
    \label{A:flamegraph}
\end{figure}
$$

$$
\begin{table}[h]
\centering
\begin{tabular}{lrrrr}
\toprule
 & est. ratio & $|C|$ & est. opt & time \\
\midrule
mean & 1.3142 & 80927.77 & 61579.82 & 168 sec\\
std & 0.0006 & 55.84 & 32.98 & 3 sec\\
min & 1.3128 & 80786 & 61511 & 161 sec\\
median & 1.3142 & 80926.50 & 61579.50 & 168 sec\\
max & 1.3157 & 81071 & 61654 & 176 sec\\
\bottomrule
\end{tabular}
\caption{Results for \textsc{Triangle Vertex Deletion} on Amazon product co-purchasing graph using rule strategy 3; fallback rule selects a random size three edge;  $n=100$}
\label{amzn_str_3_f3rand}
\end{table}
$$

$$
\begin{figure}[h]
    \centering
    \includegraphics[width=\textwidth]{./img/amazon_tvd_rules_strat.png}
    \caption{Mean number of rule executions for \textsc{Triangle Vertex Deletion} on Amazon product co-purchasing graph per rule strategy; $n=100$}
    \label{amzn_rules_strat}
\end{figure}
$$

$$
\begin{table}[h]
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.7466 & 95256 & 54539 & 3 sec\\
            std & 0.0003 & 82 & 41 & 0 sec\\
            min & 1.7458 & 95032 & 54433 & 3 sec\\
            median & 1.7466 & 95269 & 54543 & 3 sec\\
            max & 1.7474 & 95441 & 54632 & 4 sec\\
            \bottomrule
        \end{tabular}
        \caption{base rule strategy\label{amzn_str_base}}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.4731 & 86359 & 58624 & 6 sec\\
            std & 0.0005 & 61 & 36 & 0 sec\\
            min & 1.4719 & 86213 & 58534 & 6 sec\\
            median & 1.4731 & 86359 & 58627 & 6 sec\\
            max & 1.4741 & 86496 & 58723 & 6 sec\\
            \bottomrule
        \end{tabular}
        \caption{rule strategy 1\label{amzn_str_1}}
    \end{subtable}
    \newline
    \vspace{4mm}
    \newline
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.4215 & 84813 & 59666 & 38 sec\\
            std & 0.0005 & 58 & 32.33 & 0 sec\\
            min & 1.4206 & 84670 & 59586 & 37 sec\\
            median & 1.4215 & 84814 & 59664 & 38 sec\\
            max & 1.4227 & 84927 & 59740 & 39 sec\\
            \bottomrule
        \end{tabular}
        \caption{rule strategy 2\label{amzn_str_2}}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.3136 & 80829 & 61532 & 166 sec\\
            std & 0.0006 & 59 & 34.63 & 2 sec\\
            min & 1.3121 & 80684 & 61444 & 160 sec\\
            median & 1.3136 & 80836 & 61535 & 166 sec\\
            max & 1.3149 & 80943 & 61616 & 171 sec\\
            \bottomrule
        \end{tabular}
        \caption{rule strategy 3\label{amzn_str_3}}
    \end{subtable}
    \caption{Results for \textsc{Triangle Vertex Deletion} on Amazon product co-purchasing graph; $n=100$}
    \label{stats_amzn}
\end{table}
$$

\input{out/amazon_cvd_lp_naive}

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
