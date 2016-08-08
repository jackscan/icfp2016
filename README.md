### Submission of 'big.Rat' for ICFP Contest 2016

#### Source
- Used language go for developing problem solver and problem fetcher.
- Hosted at http://github.com/jackscan/icfp2016.

#### Installation
    > go get github.com/jackscan/icfp2016

#### Usage
    > icfp2016 -probdir <problems-dir> -soldir <solutions-dir>

#### Algorithm
    - add intersection points to skeleton
    - for each line s in skeleton
        - calculate transformation for s onto bottom square edge starting on left corner
        - 'a': for each possible facet in destination skeleton with line s
            - transform facet to source coordinates
            - if facet is not in unit square
                or does overlap any other polygon in source
                - continue at 'a'
            - add polygon p to source
            - for each edge e of p
                - if e is not part of any polygon in source other than p
                    and e is not on square edge
                    - recurse into 'a' with e as line s
            - if all recursions succeeded
                and all coordinates of silhouette are mapped
                - return from 'a' with success
        - if any facet lead to success
            - return solution

#### Issues
- The algorithm only considers skeleton lines and does not respect holes in the
    problem silhouette. To solve this the algorithm should also check if the
    facet chosen at 'a' overlaps with any hole polygon.
- A found solution may be too big. Algorithm should check this before returning success.
- Any bug I haven't found yet ;)
