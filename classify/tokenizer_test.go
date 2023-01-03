package main_test

import (
	"regexp"
	"testing"

	"github.com/MAA-Contest-Tester/search/classify"
)

func TestStrip1(t *testing.T) {
	source := `Let $n$ be a positive integer. Given is a subset $A$ of $\{0,1,...,5^n\}$ with $4n+2$ elements. Prove that there exist three elements $a<b<c$ from $A$ such that $c+2a>3b$. Proposed by Dominik Burek and Tomasz Ciesla, Poland`
	replace_re := regexp.MustCompile(`\s\s+`)
	tokenized := replace_re.ReplaceAllString(classify.StripSource(source), " ")
	answer := `Let n be a positive integer. Given is a subset A of 0 1 ... 5^n with 4n + 2 elements. Prove that there exist three elements a < b < c from A such that c + 2a > 3b . Proposed by Dominik Burek and Tomasz Ciesla Poland`
	if tokenized != answer {
		t.Errorf("Didn't strip correctly! Got %v", tokenized)
	}
}

func TestStrip2(t *testing.T) {
	source := `Let $I$ be the incentre of acute triangle $ABC$ with $AB\neq AC$. The incircle $\omega$ of $ABC$ is tangent to sides $BC, CA$, and $AB$ at $D, E,$ and $F$, respectively. The line through $D$ perpendicular to $EF$ meets $\omega$ at $R$. Line $AR$ meets $\omega$ again at $P$. The circumcircles of triangle $PCE$ and $PBF$ meet again at $Q$. Prove that lines $DI$ and $PQ$ meet on the line through $A$ perpendicular to $AI$. Proposed by Anant Mudgal, India`
	replace_re := regexp.MustCompile(`\s\s+`)
	tokenized := replace_re.ReplaceAllString(classify.StripSource(source), " ")
	answer := `Let  I  be the incentre of acute triangle  ABC  with  AB neq AC . The incircle   omega  of  ABC  is tangent to sides  BC  CA   and  AB  at  D  E   and  F   respectively. The line through  D  perpendicular to  EF  meets   omega  at  R . Line  AR  meets   omega  again at  P . The circumcircles of triangle  PCE  and  PBF  meet again at  Q . Prove that lines  DI  and  PQ  meet on the line through  A  perpendicular to  AI . Proposed by Anant Mudgal  India`
	answer = replace_re.ReplaceAllString(answer, " ")
	if tokenized != answer {
		t.Errorf("Didn't strip correctly! Got %v", tokenized)
	}
}
